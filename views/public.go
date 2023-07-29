package views

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed public/*
var public embed.FS

func PublicFileSystem() http.FileSystem {
	files, _ := fs.Sub(public, "public") // ignoring error!
	return http.FS(files)
}
