package http

import (
	"log"
	nethttp "net/http"

	"github.com/kaviraj-j/dispatch/internal/api/http/handlers"
	"github.com/kaviraj-j/dispatch/internal/queue"
)

type Server struct {
	addr string
	qm   *queue.Manager
	mux  *nethttp.ServeMux
}

func NewServer(addr string, qm *queue.Manager) *Server {
	s := &Server{
		addr: addr,
		qm:   qm,
		mux:  nethttp.NewServeMux(),
	}

	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	queueHandler := handlers.NewQueueHandler(s.qm)
	jobHandler := handlers.NewJobHandler(s.qm)

	// Queue routes
	s.mux.HandleFunc("/queues", queueHandler.HandleQueues)

	// Job routes
	s.mux.HandleFunc("/jobs", jobHandler.HandleJobs)
	s.mux.HandleFunc("/jobs/consume", jobHandler.HandleConsume)
	s.mux.HandleFunc("/jobs/ack", jobHandler.HandleAck)
}

func (s *Server) Start() error {
	log.Printf("HTTP server running on %s\n", s.addr)
	return nethttp.ListenAndServe(s.addr, s.mux)
}
