package main

import (
	"net/http"
	"os"
	"strconv"

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
	http.HandleFunc("GET /{$}", ui.StaticHandler.ServeHTTP)
	http.HandleFunc("POST /{$}", uploadHandler)

	err := http.ListenAndServe(listenAddress+":"+strconv.Itoa(listenPort), nil)
	if err != nil {
		logger.Error("unable to listen", "address", listenAddress, "port", listenPort, "error", err)
	}
}
