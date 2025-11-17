package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/bryan/blackjack-buddy/api/handler"
	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type server struct {
	advisor     *strategy.Advisor
	sessionStore *store.SessionStore
	ui          fs.FS
}

func newServer(strat strategy.Strategy) *server {
	return &server{
		advisor:     strategy.NewAdvisor(strat),
		sessionStore: store.NewSessionStore(),
		ui:          loadUI(),
	}
}

func (s *server) handleUI(w http.ResponseWriter, r *http.Request) {
	data, err := fs.ReadFile(s.ui, "index.html")
	if err != nil {
		http.Error(w, "trainer UI not built", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}

func (s *server) Start(port int) error {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/session", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/session" {
			handler.NewSession(s.sessionStore)(w, r)
			return
		}
		http.NotFound(w, r)
	})
	
	mux.HandleFunc("/api/session/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasSuffix(path, "/deal") && r.Method == "POST" {
			handler.NewDeal(s.sessionStore)(w, r)
			return
		}
		if strings.HasSuffix(path, "/move") && r.Method == "POST" {
			handler.NewMove(s.sessionStore)(w, r)
			return
		}
		if strings.HasSuffix(path, "/hint") && r.Method == "GET" {
			handler.NewHint(s.sessionStore, s.advisor)(w, r)
			return
		}
		http.NotFound(w, r)
	})

	fileServer := http.FileServer(http.FS(s.ui))
	mux.Handle("/assets/", fileServer)
	mux.HandleFunc("/", s.handleUI)

	addr := fmt.Sprintf(":%d", port)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Printf("Server running on http://localhost%s\n", addr)
	return srv.ListenAndServe()
}
