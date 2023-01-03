package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/techschool/simple-bank/db/sqlc"
	"github.com/techschool/simple-bank/token"
	"github.com/techschool/simple-bank/util"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     *util.Config
}

func NewServer(config *util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	server := &Server{store: store,
		tokenMaker: tokenMaker,
		router:     gin.Default(),
		config:     config,
	}

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpRouter()

	return server, nil
}

func (server *Server) setUpRouter() {
	//router := gin.Default()

	server.router.POST("/users/login", server.loginUser)
	server.router.POST("/users", server.createUser)

	authRoutes := server.router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.POST("/transfers", server.createTransfer)

	//server.router = router

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
