package main

import (
	"assignment-1/handler"
	"log"
	"net/http"
	"os"
	"time"
)

// Start time for diagnostics
var start time.Time = time.Now()

func main() {

	// Handle port assignment
	port := os.Getenv(("PORT"))
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"
	}

	// Set up handler enpoints
	http.HandleFunc(handler.DEFAULT_PATH, handler.DefaultHandler)
	http.HandleFunc(handler.DIAG_PATH, func(w http.ResponseWriter, r *http.Request) { handler.DiagHandler(w, r, start) })
	http.HandleFunc(handler.UNIINFI_PATH, handler.UniinfoHandler)
	http.HandleFunc(handler.NEIGHBOURUNIS_PATH, handler.NeighbourunisHandler)

	// Start server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
