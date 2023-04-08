package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/app"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/ids79/otus_go_homework/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	logg logger.Logg
	app  app.Application
	conf config.Config
	srv  *http.Server
}

func NewServer(logger logger.Logg, app app.Application, config config.Config) *Server {
	return &Server{
		logg: logger,
		app:  app,
		conf: config,
	}
}

func (s *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:         s.conf.HTTPServer.Address,
		Handler:      s,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.srv = server
	server.ListenAndServe()
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logg.Info("calendar is stopping...")
	s.srv.Close()
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		fmt.Println("hello-world")
		s.logg.Info("hello-world")
		s.logg.Info(ip, time.Now(), r.Method, http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("hello-world")
	}
}
