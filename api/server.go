package main

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/bryan/blackjack-buddy/api/handler"
	"github.com/bryan/blackjack-buddy/api/store"
)

type server struct {
	sessionStore *store.SessionStore
	ui           fs.FS
}

func newServer() *server {
	return &server{
		sessionStore: store.NewSessionStore(),
		ui:           loadUI(),
	}
}

func (s *server) handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/game", handler.NewGame(s.sessionStore))
	mux.HandleFunc("POST /api/game/{id}/deal", handler.NewDeal(s.sessionStore))
	mux.HandleFunc("POST /api/game/{id}/move", handler.NewMove(s.sessionStore))
	mux.HandleFunc("GET /api/game/{id}/hint", handler.NewHint(s.sessionStore))
	mux.HandleFunc("POST /api/game/{id}/abandon", handler.NewAbandon(s.sessionStore))

	fileServer := http.FileServer(http.FS(s.ui))
	mux.Handle("/assets/", fileServer)
	mux.Handle("GET /{$}", fileServer)
	return mux
}

func (s *server) Start(port int) error {
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server running on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s.handler())
}
