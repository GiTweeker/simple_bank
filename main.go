package main

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/techschool/simple-bank/api"
	db "github.com/techschool/simple-bank/db/sqlc"
	_ "github.com/techschool/simple-bank/docs/statik"
	"github.com/techschool/simple-bank/gapi"
	"github.com/techschool/simple-bank/pb"
	"github.com/techschool/simple-bank/util"
	"github.com/techschool/simple-bank/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {

	config, err := util.LoadConfig(".")
	if strings.ToLower(config.Environment) == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if err != nil {
		log.Fatal().Msg("Cannot load config " + err.Error())
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to db: " + err.Error())
	}

	runDbMigration(config.MigrationPath, config.DBSource)

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisServer,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcess(redisOpt, store)
	//runGinServer(config, store)
	go runGrpcGatewayServer(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)

}

func runDbMigration(path string, source string) {
	migration, err := migrate.New(path, source)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migration instance" + err.Error())
	}

	if err = migration.Up(); err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("cannot run new migration up " + err.Error())
	}
	log.Info().Msg("db migrated successfully")
}
func runGrpcGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(&config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("Cannot start server: " + err.Error())
	}
	gatewayMarshallOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(
		gatewayMarshallOptions,
	)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start  handler server: " + err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	//fs := http.FileServer(http.Dir("./docs/swagger"))

	statickFs, err := fs.New()

	if err != nil {
		log.Fatal().Msg("Cannot load static files  : " + err.Error())
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statickFs))

	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HttpServerAddress)

	if err != nil {
		log.Fatal().Msg("Cannot create listener : " + err.Error())
	}

	log.Info().Msgf("Started gRPC Gateway Server at %s", listener.Addr().String())

	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create GRPC Gateway Server")
	}

}
func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(&config, store, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server: " + err.Error())
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterSimpleBankServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create listener : " + err.Error())
	}

	log.Printf("Started gRPC Server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create GRPC Server")
	}

}
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(&config, store)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server: " + err.Error())
	}
	err = server.Start(config.HttpServerAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server: " + err.Error())
	}
}

func runTaskProcess(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task processor")
	}

}
