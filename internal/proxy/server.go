package proxy

import (
	"github.com/BroNaz/proxy/internal/config"
	"net/http"
)

type Server struct {
	httpClient *http.Client
	config config.Settings
}