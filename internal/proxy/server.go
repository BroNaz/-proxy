package proxy

import (
	"context"
	"net/http"
	"time"

	"github.com/BroNaz/proxy/internal/config"
	"github.com/rs/zerolog/log"
)

type Server struct {
	httpClient *http.Client
	config     config.Settings
}

func NewServer(config config.Settings) (*Server, error) {
	log.Debug().Msg("Server began to form")

	server := new(Server)
	server.config = config

	server.httpClient = new(http.Client)
	server.httpClient.Timeout = 5 * time.Second

	log.Debug().Msg("Server formed")
	return server, nil
}

func (s *Server) Run() error {
	log.Info().Msg("Server run")

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("gently stopping the server")
	return nil
}
