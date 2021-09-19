package proxy

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/BroNaz/proxy/internal/config"
	"github.com/rs/zerolog/log"
)

type Server struct {
	httpServer *http.Server
	config     config.Settings
}

func NewServer(config config.Settings, addr string) *Server {
	log.Debug().Msg("Server began to form")

	s := new(Server)
	s.config = config

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodConnect {
			handleTunneling(w, r)
		} else {
			handleHTTP(w, r)
		}
	})

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Debug().Msg("Server formed")
	return s
}

func (s *Server) Run() error {
	go s.httpServer.ListenAndServe()

	var err error
	if !s.config.HTTPS {
		err = s.httpServer.ListenAndServe()
	} else {
		err = s.httpServer.ListenAndServeTLS("ssl/server.pem", "ssl/server.key")
	}

	if err != nil {
		log.Warn().
			Str("error", err.Error()).
			Msg("the server did not start")
		return err
	}

	log.Info().
		Str("addr", s.httpServer.Addr).
		Msg("server starts listening")

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("gently stopping the server")
	return nil
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
