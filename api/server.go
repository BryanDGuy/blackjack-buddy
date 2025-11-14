package main

import (
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"
	"time"

	"github.com/bryan/blackjack-buddy/api/handler"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type server struct {
	advisor *strategy.Advisor
	engine  *game.Engine
	ui      fs.FS
}

func newServer(strat strategy.Strategy) *server {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &server{
		advisor: strategy.NewAdvisor(strat),
		engine:  game.NewEngine(rng),
		ui:      loadUI(),
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
	mux.HandleFunc("/api/deal", handler.NewDeal(s.engine, s.advisor))
	mux.HandleFunc("/api/check", handler.NewCheck(s.advisor, s.engine))
	mux.HandleFunc("/api/hint", handler.NewHint(s.advisor))

	fileServer := http.FileServer(http.FS(s.ui))
	mux.Handle("/assets/", fileServer)
	mux.HandleFunc("/", s.handleUI)

	addr := fmt.Sprintf(":%d", port)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return srv.ListenAndServe()
}
