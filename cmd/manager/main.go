package main

import (
	"SIEMSEasy/internal/database"
	"SIEMSEasy/internal/realtime"
	"context"
	"crypto/sha256"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
)

func main() {

	ctx := context.Background()

	testUUID, _ := uuid.NewRandom()
	testHash := sha256.Sum256([]byte("Hello world"))
	testIP := net.IP([]byte{0xff, 0xff, 0xff, 0xff})
	testLog := database.LogSchema{
		LogID:           testUUID,
		DeviceID:        testUUID,
		SourceAddress:   testIP,
		ProductName:     []byte("linux"),
		ServiceName:     []byte("sshd"),
		Message:         []byte("testing"),
		MessageHash:     testHash,
		ScrapeTimestamp: time.Now(),
		RecvTimestamp:   time.Now(),
	}

	mainBuffer := realtime.CreateRingBuffer()
	mainBuffer.Enqueue(testLog)

	conn, err := database.DbConnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	database.DbInit(ctx, conn)

}
