package api

import (
	db "github.com/GritNg/simple-bank-backend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store    // for interact with database
	router *gin.Engine // for routing
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/account/create", server.createAccount)
	router.GET("/account/:id", server.getAccount)
	router.GET("/account/list", server.listAccount)
	router.POST("/transfer", server.createTransfer)

	router.POST("/user/create", server.createUser)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
