package main

import (
	"SIEMSEasy/internal/database"
	"SIEMSEasy/internal/network/tcpserver"
	"SIEMSEasy/internal/realtime"
	"context"
	"crypto/sha256"
	"fmt"
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

	fmt.Println("Initialising Ring Buffer")
	mainBuffer := realtime.CreateRingBuffer()
	mainBuffer.Enqueue(testLog)
	fmt.Println("Ring Buffer Ready")

	fmt.Println("Initialising Database Connection...")
	conn, err := database.DbConnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("Successful ClickHouse Connection")

	database.DbInit(ctx, conn)

	tcpserver.ListenIncoming()

}
