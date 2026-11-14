package http_server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type HealthCheck struct {
	Status    bool   `json:"status"`
	Timestamp string `json:"timestamp"`
}

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
