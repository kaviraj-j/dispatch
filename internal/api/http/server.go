package http

import (
	"log"
	nethttp "net/http"

	"github.com/kaviraj-j/dispatch/internal/api/http/handlers"
	"github.com/kaviraj-j/dispatch/internal/auth"
	"github.com/kaviraj-j/dispatch/internal/config"
	"github.com/kaviraj-j/dispatch/internal/middleware"
	"github.com/kaviraj-j/dispatch/internal/queue"
)

type Server struct {
	addr       string
	qm         *queue.Manager
	mux        *nethttp.ServeMux
	middleware *middleware.Middleware
}

func NewServer(config config.Config, qm *queue.Manager) *Server {
	addr := config.ServerAddress
	auth := auth.NewAuth(config.ProducerApiKey, config.ConsumerApiKey)
	s := &Server{
		addr:       addr,
		qm:         qm,
		mux:        nethttp.NewServeMux(),
		middleware: middleware.NewMiddleware(auth),
	}

	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	queueHandler := handlers.NewQueueHandler(s.qm)
	jobHandler := handlers.NewJobHandler(s.qm)

	s.mux.Handle(
		"/queues",
		s.middleware.IsProducerAuthenticated(
			nethttp.HandlerFunc(queueHandler.HandleQueues),
		),
	)

	s.mux.Handle(
		"/jobs",
		s.middleware.IsProducerAuthenticated(
			nethttp.HandlerFunc(jobHandler.HandleJobs),
		),
	)

	s.mux.Handle(
		"/jobs/consume",
		s.middleware.IsConsumerAuthenticated(
			nethttp.HandlerFunc(jobHandler.HandleConsume),
		),
	)

	s.mux.Handle(
		"/jobs/ack",
		s.middleware.IsConsumerAuthenticated(
			nethttp.HandlerFunc(jobHandler.HandleAck),
		),
	)
}

func (s *Server) Start() error {
	log.Printf("HTTP server running on %s\n", s.addr)
	return nethttp.ListenAndServe(s.addr, s.mux)
}
