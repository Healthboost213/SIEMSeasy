package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

/*

SETUP ROUTES FOR:
- LOG INGEST
- AUTHENTICATE
-

*/

// JSON response defs

type HealthCheck struct {
	Status    bool   `json:"status"`
	Timestamp string `json:"timestamp"`
}

// Manager Routes

func healthCheck(w http.ResponseWriter, r *http.Request) {

	currHealth := HealthCheck{
		Status:    true,
		Timestamp: time.Now().Format(time.RFC3339Nano),
	}

	bodyData, err := json.Marshal(currHealth)
	if err != nil {
		log.Fatal("JSON Error:", err)
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bodyData)

}

// Worker Routes
