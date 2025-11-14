package main

import (
	"flag"
	"log"

	"github.com/bryan/blackjack-buddy/internal/strategy/strategies"
)

func main() {
	port := flag.Int("port", 8080, "Port for the HTTP server")
	strategyName := flag.String("strategy", string(strategies.BasicStrategy), "Strategy variant")
	flag.Parse()

	strat := strategies.CreateStrategy(strategies.StrategyType(*strategyName))
	srv := newServer(strat)

	if err := srv.Start(*port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
