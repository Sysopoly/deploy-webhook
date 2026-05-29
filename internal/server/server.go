package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sysopoly/deploy-webhook/internal/config"
	"github.com/sysopoly/deploy-webhook/internal/runner"
	"github.com/sysopoly/deploy-webhook/internal/webhook"
)

type Server struct {
	cfg    *config.Config
	runner *runner.Runner
	mux    *http.ServeMux
}

func New(cfg *config.Config) *http.Server {
	s := &Server{
		cfg:    cfg,
		runner: runner.New(cfg.ApplyScript),
		mux:    http.NewServeMux(),
	}

	s.mux.HandleFunc("POST /webhook", s.handleWebhook)
	s.mux.HandleFunc("GET /health", s.handleHealth)

	return &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: s.mux,
	}
}

func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := webhook.ValidateSignature(r, s.cfg.Secret)
	if err != nil {
		log.Printf("Signature validation failed: %v", err)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var payload struct {
		Ref string `json:"ref"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Bad payload", http.StatusBadRequest)
		return
	}

	expectedRef := fmt.Sprintf("refs/heads/%s", s.cfg.Branch)
	if payload.Ref != expectedRef {
		log.Printf("Ignoring push to %s (want %s)", payload.Ref, expectedRef)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ignored","reason":"wrong branch"}`)
		return
	}

	s.runner.Trigger()
	log.Printf("Triggered apply for push to %s", s.cfg.Branch)

	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, `{"status":"accepted"}`)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	status := s.runner.Status()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"healthy","apply_running":%t,"pending":%t}`, status.Running, status.Pending)
}
