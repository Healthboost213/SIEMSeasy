package httpserver

import (
	"encoding/json"
	"fmt"
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

type LogJSON struct {
	LogID           []byte `json:"logID"`
	DeviceID        []byte `json:"deviceID"`
	SourceAddress   []byte `json:"sourceAddress"`
	ProductName     string `json:"productName"`
	ServiceName     string `json:"serviceName"`
	Message         string `json:"message"`
	MessageHash     []byte `json:"messageHash"`
	ScrapeTimestamp string `json:"scrapeTimestamp"`
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

func ingestLogs(w http.ResponseWriter, r *http.Request) {

	var message LogJSON

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)
	if err != nil {
		log.Fatal("JSON Decode Error:", err)
	}
	fmt.Println(message.LogID)

}

// Worker Routes
