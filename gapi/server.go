package gapi

import (
	"fmt"
	db "github.com/techschool/simple-bank/db/sqlc"
	"github.com/techschool/simple-bank/pb"
	"github.com/techschool/simple-bank/token"
	"github.com/techschool/simple-bank/util"
	"github.com/techschool/simple-bank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store           db.Store
	tokenMaker      token.Maker
	config          *util.Config
	taskDistributor worker.TaskDistributor
}

func NewServer(config *util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:           store,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
