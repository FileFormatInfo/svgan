package main

import (
	"net/http"

	"github.com/FileFormatInfo/svgan/internal/common"
	svgan "github.com/FileFormatInfo/svgan/lib"
	"github.com/FileFormatInfo/svgan/ui"
)

func uploadHandler(w http.ResponseWriter, r *http.Request, e error) {
	ui.RunTemplate(w, r, "upload.tmpl", map[string]interface{}{
		"Err":   e,
		"Title": "SVG File Upload",
	})
}

func uploadGetHandler(w http.ResponseWriter, r *http.Request) {
	uploadHandler(w, r, nil)
}

func uploadPostHandler(w http.ResponseWriter, r *http.Request) {

	// max = 10 MB
	parseErr := r.ParseMultipartForm(10 << 20)
	if parseErr != nil {
		uploadHandler(w, r, parseErr)
		return
	}

	file, handler, fileErr := r.FormFile("file")
	if fileErr != nil {
		uploadHandler(w, r, fileErr)
		return
	}
	defer file.Close()

	// seems wasteful, but handler.content is private...
	raw := make([]byte, handler.Size)
	_, readErr := file.Read(raw)
	if readErr != nil {
		uploadHandler(w, r, readErr)
		return
	}

	svgInfo, svgErr := svgan.SvgCheck(common.Logger, raw)
	if svgErr != nil {
		uploadHandler(w, r, svgErr)
		return
	}

	ui.RunTemplate(w, r, "_results.tmpl", map[string]interface{}{
		"filename": handler.Filename,
		"size":     handler.Size,
		"mime":     handler.Header.Get("Content-Type"),
		"data":     svgInfo,
		"Title":    "SVG Analysis Results",
	})
}
