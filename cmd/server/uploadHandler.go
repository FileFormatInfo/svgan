package main

import (
	"fmt"
	"net/http"

	svgan "github.com/FileFormatInfo/svgan/lib"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	// max = 10 MB
	parseErr := r.ParseMultipartForm(10 << 20)
	if parseErr != nil {
		http.Error(w, fmt.Sprintf("Unable to parse form data (%v)\n", parseErr), http.StatusBadRequest)
		return
	}

	file, handler, fileErr := r.FormFile("file")
	if fileErr != nil {
		http.Error(w, fmt.Sprintf("Unable to get file from form (%v)\n", fileErr), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	defer file.Close()

	fmt.Fprintf(w, "File: %+v\n", handler.Filename)
	fmt.Fprintf(w, "Size: %+v\n", handler.Size)
	fmt.Fprintf(w, "MIME: %+v\n", handler.Header.Get("Content-Type"))

	// seems wasteful, but handler.content is private...
	raw := make([]byte, handler.Size)
	_, readErr := file.Read(raw)
	if readErr != nil {
		fmt.Fprintf(w, "ERROR: Reading the File (%v)\n", readErr)
		return
	}

	svgInfo, svgErr := svgan.SvgCheck(logger, raw)
	if svgErr != nil {
		fmt.Fprintf(w, "ERROR: Parsing the SVG (%v)\n", svgErr)
		return
	}

	fmt.Fprintf(w, "INFO: %+v\n", svgInfo)
}
