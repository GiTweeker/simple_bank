package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/techschool/simple-bank/api"
	db "github.com/techschool/simple-bank/db/sqlc"
	"github.com/techschool/simple-bank/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot load config ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	server, err := api.NewServer(&config, store)

	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}
	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server: ", err)
	}

}
