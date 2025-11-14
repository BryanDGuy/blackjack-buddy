package web

import (
	"embed"
	"io/fs"
)

var (
	//go:embed dist/* dist/assets/*
	dist embed.FS
)

func Dist() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		return nil
	}
	return sub
}
