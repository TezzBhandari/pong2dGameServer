package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/TezzBhandari/pong"
	"github.com/gorilla/mux"
)

const shutdownTimeout = 10 * time.Second

type Server struct {
	server *http.Server
	router *mux.Router
	ln     net.Listener
	addr   string
}

func NewHttpServer(addr string) *Server {

	s := &Server{
		server: &http.Server{
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 1 * time.Second,
			IdleTimeout:  10 * time.Second,
		},
		router: mux.NewRouter(),
		addr:   addr,
	}

	s.router.Use(reportPanic)
	s.server.Handler = s.ServeHttp()
	s.router.NotFoundHandler = s.handleNotFound()
	router := s.router.PathPrefix("/").Subrouter()
	router.Handle("/", s.handleIndex())

	return s
}

func (s *Server) Open() error {
	var err error

	if s.ln, err = net.Listen("tcp", s.addr); err != nil {
		return err
	}
	// Begin serving requests on the listener. We use Serve() instead of
	// ListenAndServe() because it allows us to check for listen errors (such
	// as trying to use an already open port) synchronously.
	go s.server.Serve(s.ln)
	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) ServeHttp() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		s.router.ServeHTTP(rw, r)
	})
}

func reportPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// http.Error(rw, "Internal Error", http.StatusInternalServerError)
				rw.WriteHeader(http.StatusInternalServerError)
				pong.ReportPanic(err)
			}
		}()
		next.ServeHTTP(rw, r)
	})
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (s *Server) handleNotFound() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(rw).Encode(ErrorResponse{Message: "route not found"})
		if err != nil {
			fmt.Printf("[ResponseError] %q", err)
		}
	})
}

type SuccessResponse struct {
	Data map[string]any `json:"data"`
}

func (s *Server) handleIndex() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		err := json.NewEncoder(rw).Encode(SuccessResponse{
			Data: map[string]any{
				"status": "ok",
			},
		})

		if err != nil {
			fmt.Printf("[ResponseError] %q", err)
		}
	})
}

func (s *Server) Port() int {
	if s.ln != nil {
		return 0

	}
	return s.ln.Addr().(*net.TCPAddr).Port
}

func (s *Server) URL() string {
	return fmt.Sprintf("http://localhost:%d", s.Port())
}
