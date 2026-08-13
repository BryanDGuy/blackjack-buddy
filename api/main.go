package main

import (
	"flag"
	"log"
)

func main() {
	port := flag.Int("port", 8080, "Port for the HTTP server")
	flag.Parse()

	if err := newServer().Start(*port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
