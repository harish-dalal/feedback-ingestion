package server

import (
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Server struct {
	Router *http.ServeMux
	DBPool *pgxpool.Pool
}

func NewServer(dbpool *pgxpool.Pool) *Server {
	return &Server{
		Router: http.NewServeMux(),
		DBPool: dbpool,
	}
}
