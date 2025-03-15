package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/FileFormatInfo/svgan/internal/common"
	svgan "github.com/FileFormatInfo/svgan/lib"
	"github.com/FileFormatInfo/svgan/ui"
)

func urlHandler(w http.ResponseWriter, r *http.Request, e error) {
	ui.RunTemplate(w, r, "url.tmpl", map[string]interface{}{
		"Err":   e,
		"Title": "URL Analysis",
	})
}

func urlGetHandler(w http.ResponseWriter, r *http.Request) {
	urlHandler(w, r, nil)
}

func urlPostHandler(w http.ResponseWriter, r *http.Request) {

	url := r.FormValue("url")
	if url == "" {
		urlHandler(w, r, errors.New("no URL provided"))
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		urlHandler(w, r, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		urlHandler(w, r, errors.New("failed to fetch URL"))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		urlHandler(w, r, err)
		return
	}

	svgInfo, svgErr := svgan.SvgCheck(common.Logger, body)
	if svgErr != nil {
		urlHandler(w, r, svgErr)
		return
	}

	ui.RunTemplate(w, r, "_results.tmpl", map[string]interface{}{
		"source": url,
		"size":   len(body),
		"mime":   resp.Header.Get("Content-Type"),
		"data":   svgInfo,
		"Title":  "URL Analysis Results",
	})
}

func urlJsonHandler(w http.ResponseWriter, r *http.Request) {
	var url string

	if r.Method == http.MethodPost {
		url = r.FormValue("url")
	} else {
		url = r.URL.Query().Get("url")
	}
	if url == "" {
		handleJson(w, r, map[string]interface{}{
			"success": false,
			"error":   "url parameter missing",
		})
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		handleJson(w, r, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"source":  url,
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		handleJson(w, r, map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("failed to fetch URL (http response=%d)", resp.StatusCode),
			"source":  url,
		})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		handleJson(w, r, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"source":  url,
		})
		return
	}

	svgInfo, svgErr := svgan.SvgCheck(common.Logger, body)
	if svgErr != nil {
		handleJson(w, r, map[string]interface{}{
			"success": false,
			"error":   svgErr.Error(),
			"source":  url,
		})
		return
	}

	handleJson(w, r, map[string]interface{}{
		"success": true,
		"source":  url,
		"size":    len(body),
		"mime":    resp.Header.Get("Content-Type"),
		"data":    svgInfo,
	})
}
