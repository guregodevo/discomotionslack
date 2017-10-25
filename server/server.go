package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/guregodevo/discomotionslack/misc"
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

type ServerInfo struct {
	Server    string `json:"server"`
	Version   string `json:"version"`
	Build     string `json:"build"`
	BuildTime string `json:"buildTime"`
	Hostname  string `json:"hostname"`
}

type Server struct {
	Info    *ServerInfo
	BaseURL string
	Uptime  time.Time
	Server  *http.Server
	Api     *slack.Client
	Mutex   *sync.Mutex
}

func (s *Server) Run(conf *misc.Http) error {
	router := s.NewRouter()

	s.Server = &http.Server{
		Addr:              conf.Address,
		ReadTimeout:       time.Duration(conf.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(conf.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(conf.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(conf.ReadTimeout*2) * time.Second,
		Handler:           router,
	}

	log.WithField("address", conf.Address[1:]).Info("HTTP server ready")

	return s.Server.ListenAndServe()
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if s.Server != nil {
		if err := s.Server.Shutdown(ctx); err != nil {
			log.WithField("error", err).Error("Failed to shut down http server cleanly")

			// Close all open connections
			if err := s.Server.Close(); err != nil {
				log.WithField("error", err).Error("Failed to force-close http server")
			}
		}
	}

	return nil
}

func (s *Server) NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	s.addRoute(router, "slackPlayInteractive", "POST", "/discomotion/v1/interactive", s.Interactive)
	s.addRoute(router, "slackPlay", "POST", "/discomotion/v1/play", s.Play)
	s.addRoute(router, "Index", "GET", "/", s.Index)

	return router
}

func (s *Server) addRoute(router *mux.Router, name string, method string, pattern string, fn http.HandlerFunc) {
	var handler http.Handler
	handler = s.WrapRequest(fn, name)

	router.Methods(method).Path(pattern).Name(name).Handler(handler)

}

type LoggedWriter struct {
	http.ResponseWriter

	statusCode int
}

func NewLoggedWriter(w http.ResponseWriter) *LoggedWriter {
	return &LoggedWriter{w, http.StatusOK}
}

func (w *LoggedWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode

	w.ResponseWriter.WriteHeader(statusCode)
}

func (s *Server) WrapRequest(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := NewLoggedWriter(w)
		inner.ServeHTTP(lw, r)

		duration := time.Since(start)
		statusCode := lw.statusCode

		log.WithFields(log.Fields{
			"ip":       misc.GetIPAdress(r),
			"method":   r.Method,
			"uri":      r.RequestURI,
			"name":     name,
			"status":   statusCode,
			"duration": duration,
		}).Info(name)

	})
}
