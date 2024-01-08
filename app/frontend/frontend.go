package frontend

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var frontend embed.FS

type Frontend http.FileSystem

func NewFilesystemBackedFrontend() Frontend {
	return http.Dir("app/frontend/static")
}
func NewEmbedBackendFrontend() Frontend {
	fsys, err := fs.Sub(frontend, "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
