package main

import (
	"fmt"

	"github.com/kaviraj-j/dispatch/internal/api/http"
	"github.com/kaviraj-j/dispatch/internal/config"
	"github.com/kaviraj-j/dispatch/internal/queue"
)

func main() {
	cfg := config.Load()
	qm := queue.NewManager()
	server := http.NewServer(cfg, qm)
	server.Start()
	fmt.Println("Dispatch running...")
}
