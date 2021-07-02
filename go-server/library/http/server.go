package http

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	*http.Server
}

type ServerConf struct {
	Addr string
}

func NewServer(cf *ServerConf) (svr *Server) {
	return &Server{
		Server: &http.Server{
			Addr: cf.Addr,
		},
	}
}

func (svr *Server) Run(handler http.Handler) (err error) {
	svr.Handler = handler

	err = svr.Server.ListenAndServe()

	if err == http.ErrServerClosed {
		err = nil
	}

	if err != nil {
		err = fmt.Errorf("http.Server.ListenAndServe: %w", err)
	}

	return
}

func (svr *Server) Close() (err error) {
	if err = svr.Shutdown(context.TODO()); err != nil {
		err = fmt.Errorf("server shutdown: %w", err)
		return
	}

	return
}
