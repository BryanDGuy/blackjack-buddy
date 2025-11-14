package main

import (
	"flag"
	"log"

	"github.com/bryan/blackjack-buddy/internal/strategy/strategies"
)

func main() {
	port := flag.Int("port", 8080, "Port for the HTTP server")
	flag.Parse()

	strat := strategies.NewBasic()
	srv := newServer(strat)

	if err := srv.Start(*port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
