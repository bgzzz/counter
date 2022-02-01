package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sync"

	"github.com/bgzzz/counter/pkg/model"

	log "github.com/sirupsen/logrus"
)

const (
	apiPrefix = "/api"

	MaxUint64 = ^uint64(0)
)

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
	port    int
	httpSrv *http.Server

	counter uint64
	mtx     sync.RWMutex
}

func (srv *Server) setupHandlers() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc(path.Join(apiPrefix, model.APIVersion, "counter"),
		srv.handleCounter)
	return mux
}

func (srv *Server) handleCounter(w http.ResponseWriter,
	r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			srv.getCounter(w)
		}
	case http.MethodPost:
		{
			srv.incrCounter(w)
		}
	case http.MethodDelete:
		{
			srv.decrCounter(w)
		}
	default:
		http.Error(w,
			fmt.Sprintf("method %s is not supported", r.Method),
			http.StatusMethodNotAllowed)
	}
}

func (srv *Server) getCounter(w http.ResponseWriter) {

	srv.mtx.RLock()
	defer srv.mtx.RUnlock()

	counter := model.CounterRsp{
		Counter: srv.counter,
	}

	log.Debugf("get the counter with value %d was executed",
		counter.Counter)

	sendCntRsp(w, &counter, http.StatusOK)
}

func (srv *Server) incrCounter(w http.ResponseWriter) {

	srv.mtx.Lock()
	defer srv.mtx.Unlock()

	if srv.counter == MaxUint64 {
		log.Debug("increment method on maximum unit64")
		http.Error(w,
			"unable to increment, counter has reached its maximum value",
			http.StatusUnprocessableEntity)
		return
	}

	srv.counter++

	log.Debugf("increment counter with value %d was executed",
		srv.counter)

	sendCntRsp(w, &model.CounterRsp{Counter: srv.counter},
		http.StatusCreated)
}

func (srv *Server) decrCounter(w http.ResponseWriter) {

	srv.mtx.Lock()
	defer srv.mtx.Unlock()

	if srv.counter == 0 {
		log.Debug("decrement method on 0")
		http.Error(w,
			"unable to decrement, counter has reached its minimum value",
			http.StatusUnprocessableEntity)
		return
	}

	srv.counter--

	log.Debugf("decrement counter with value %d was executed",
		srv.counter)

	sendCntRsp(w, &model.CounterRsp{Counter: srv.counter},
		http.StatusOK)
}

func (srv *Server) Run() error {
	log.Infof("Running server on port: %d", srv.port)
	return srv.httpSrv.ListenAndServe()
}

func sendCntRsp(w http.ResponseWriter,
	counter *model.CounterRsp, status int) {
	b, err := json.Marshal(counter)
	if err != nil {
		log.Errorf("unable to marshal: %v", err)
		http.Error(w,
			"something went wrong",
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(b); err != nil {
		log.Errorf("unable to write response: %v", err)
		return
	}
	log.Debugf("get the counter with value %d was executed",
		counter.Counter)
}
