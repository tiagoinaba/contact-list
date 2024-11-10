package server

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/tiagoinaba/davinti/internal/database"
)

type Server struct {
	Router *gin.Engine
	DB     *sql.DB
}

func (s *Server) Run() error {
	return s.Router.Run()
}

func New() *Server {
	return &Server{
		Router: gin.Default(),
		DB:     database.New(),
	}
}
