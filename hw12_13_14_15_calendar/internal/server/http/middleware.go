package internalhttp

import (
	"net"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		s.logg.Info(ip, time.Now(), " ", r.Method, " ", http.StatusOK)
		s.ServeHTTP(w, r)
	})
}
