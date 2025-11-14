package main

import (
	"embed"
	"io/fs"
)

var (
	//go:embed assets/* assets/assets/*
	ui embed.FS
)

func loadUI() fs.FS {
	sub, err := fs.Sub(ui, "assets")
	if err != nil {
		panic("ui assets missing: run make ui-build before starting the server")
	}
	return sub
}
