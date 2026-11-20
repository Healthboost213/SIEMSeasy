package httpserver

import (
	"log"
	"net/http"
)

func StartManagerHTTP() {

	mux := http.NewServeMux()

	// Manager Specific Routes
	mux.HandleFunc("GET /api/health", healthCheck)

	// Manager takes port 6250
	if err := http.ListenAndServe(":6210", mux); err != nil {
		log.Println("HTTP Server Error:", err)
	}

}
