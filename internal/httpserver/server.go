package httpserver

import (
	"log"
	"net/http"
)

func StartManagerHTTP() {

	mux := http.NewServeMux()

	// Manager Specific Routes
	mux.HandleFunc("GET /api/health", healthCheck)
	mux.HandleFunc("POST /api/ingest", ingestLogs)

	// Manager takes port 6250
	if err := http.ListenAndServe(":6250", mux); err != nil {
		log.Fatal("HTTP Server Error:", err)
	}

}
