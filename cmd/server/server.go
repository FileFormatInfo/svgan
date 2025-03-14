package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/FileFormatInfo/svgan/internal/common"
	"github.com/FileFormatInfo/svgan/ui"
)

func main() {

	var listenPort, portErr = strconv.Atoi(os.Getenv("PORT"))
	if portErr != nil {
		listenPort = 4000
	}
	var listenAddress = os.Getenv("ADDRESS")

	http.HandleFunc("/status.json", statusHandler)
	http.HandleFunc("/robots.txt", ui.StaticHandler.ServeHTTP)
	http.HandleFunc("/favicon.ico", ui.StaticHandler.ServeHTTP)
	http.HandleFunc("/favicon.svg", ui.StaticHandler.ServeHTTP)
	http.HandleFunc("/images/", ui.StaticHandler.ServeHTTP)
	http.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/index.html", http.StatusSeeOther) })
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) { ui.RunTemplate(w, r, "index.tmpl", nil) })
	http.HandleFunc("GET /upload.html", uploadGetHandler)
	http.HandleFunc("POST /upload.html", uploadPostHandler)
	http.HandleFunc("GET /url.html", urlGetHandler)
	http.HandleFunc("POST /url.html", urlPostHandler)
	http.HandleFunc("GET /clipboard.html", clipboardGetHandler)
	http.HandleFunc("POST /clipboard.html", clipboardPostHandler)

	err := http.ListenAndServe(listenAddress+":"+strconv.Itoa(listenPort), nil)
	if err != nil {
		common.Logger.Error("unable to listen", "address", listenAddress, "port", listenPort, "error", err)
	}
}
