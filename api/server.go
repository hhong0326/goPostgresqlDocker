package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/hhong0326/goPostgresqlDocker.git/db/sqlc"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	// add routes to router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.PUT("/accounts", server.updateAccount)
	router.DELETE("/accounts", server.deleteAccount)
	server.router = router

	return server
}

// Start runs the HTTP server on a speciifc address
func (server *Server) Start(address string) error {

	return server.router.Run(address)
}

// errorResponse deals with errors
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
