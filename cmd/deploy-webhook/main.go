package main

import (
	"log"
	"os"

	"github.com/sysopoly/deploy-webhook/internal/config"
	"github.com/sysopoly/deploy-webhook/internal/server"
)

func main() {
	cfg := config.Load()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
	log.SetPrefix("[deploy-webhook] ")

	srv := server.New(cfg)
	log.Printf("Starting on :%s (branch filter: %s)", cfg.Port, cfg.Branch)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
