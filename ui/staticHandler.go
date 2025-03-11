package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var embeddedFiles embed.FS
var StaticHandler http.Handler

func initStaticHandler() (http.Handler, error) {

	fsys, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		return nil, err
	}

	return http.FileServer(http.FS(fsys)), nil
}

func init() {
	var err error
	StaticHandler, err = initStaticHandler()
	if err != nil {
		panic(err)
	}
}
