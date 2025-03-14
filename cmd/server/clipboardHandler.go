package main

import (
	"errors"
	"net/http"

	"github.com/FileFormatInfo/svgan/internal/common"
	svgan "github.com/FileFormatInfo/svgan/lib"
	"github.com/FileFormatInfo/svgan/ui"
)

func clipboardHandler(w http.ResponseWriter, r *http.Request, e error) {
	ui.RunTemplate(w, r, "clipboard.tmpl", map[string]interface{}{
		"Err":   e,
		"Title": "SVG Text Analysis",
	})
}

func clipboardGetHandler(w http.ResponseWriter, r *http.Request) {
	clipboardHandler(w, r, nil)
}

func clipboardPostHandler(w http.ResponseWriter, r *http.Request) {
	clipboard := r.FormValue("xml")
	if clipboard == "" {
		clipboardHandler(w, r, errors.New("no clipboard data provided"))
		return
	}

	svgInfo, svgErr := svgan.SvgCheck(common.Logger, []byte(clipboard))
	if svgErr != nil {
		clipboardHandler(w, r, svgErr)
		return
	}

	ui.RunTemplate(w, r, "_results.tmpl", map[string]interface{}{
		"source": "clipboard",
		"size":   len(clipboard),
		"data":   svgInfo,
		"Title":  "SVG Analysis Results",
	})
}
