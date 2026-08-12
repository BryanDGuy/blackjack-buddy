package main

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/bryan/blackjack-buddy/api/handler"
	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type server struct {
	advisor      *strategy.Advisor
	sessionStore *store.SessionStore
	ui           fs.FS
}

func newServer(strat strategy.Strategy) *server {
	return &server{
		advisor:      strategy.NewAdvisor(strat),
		sessionStore: store.NewSessionStore(),
		ui:           loadUI(),
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

func (s *server) handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/game", handler.NewGame(s.sessionStore))
	mux.HandleFunc("POST /api/game/{id}/deal", handler.NewDeal(s.sessionStore))
	mux.HandleFunc("POST /api/game/{id}/move", handler.NewMove(s.sessionStore))
	mux.HandleFunc("GET /api/game/{id}/hint", handler.NewHint(s.sessionStore, s.advisor))
	mux.HandleFunc("POST /api/game/{id}/abandon", handler.NewAbandon(s.sessionStore))

	fileServer := http.FileServer(http.FS(s.ui))
	mux.Handle("/assets/", fileServer)
	mux.HandleFunc("GET /{$}", s.handleUI)
	return mux
}

func (s *server) Start(port int) error {

	addr := fmt.Sprintf(":%d", port)

	srv := &http.Server{
		Addr:    addr,
		Handler: s.handler(),
	}

	fmt.Printf("Server running on http://localhost%s\n", addr)
	return srv.ListenAndServe()
}
