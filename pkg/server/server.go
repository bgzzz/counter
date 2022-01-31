package server

import (
	"fmt"
	"net/http"
	"path"

	log "github.com/sirupsen/logrus"
)

const apiPrefix = "/api/v1"

func NewServer(port int) *Server {

	srv := Server{
		port: port,
	}
	srv.httpSrv = &http.Server{
		Addr:    fmt.Sprintf(":%d", srv.port),
		Handler: srv.setupHandlers(),
	}

	return &srv
}

type Server struct {
	port int

	httpSrv *http.Server
}

func (srv *Server) setupHandlers() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc(path.Join(apiPrefix, "counter"),
		srv.handleCounter)
	return mux
}

func (srv *Server) handleCounter(w http.ResponseWriter,
	r *http.Request) {
	fmt.Fprintf(w, "counter")
}

func (srv *Server) Run() error {
	log.Infof("Running server on port: %d", srv.port)
	return srv.httpSrv.ListenAndServe()
}
