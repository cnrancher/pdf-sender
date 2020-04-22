package server

import (
	"fmt"
	"net/http"

	"github.com/rancher/pdf-sender/pkg/apis"

	"github.com/sirupsen/logrus"
)

var limiter = NewIPRateLimiter(0.1, 4)

type Server struct {
	port int
}

func New(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Run() error {
	apis.ConnectMysql()
	apis.CreateCacheClient()
	container := apis.RegisterAPIs()
	apis.CollectInformation()

	logrus.Infof("server running, listening at :%d\n", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), limitMiddleware(container.ServeMux))
}

func limitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := limiter.GetLimiter(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
