package main

import (
	"SIEMSEasy/internal/database"
	"SIEMSEasy/internal/network/tcpserver"
	"crypto/sha256"
	"net"
	"time"

	"github.com/google/uuid"
)

func main() {

	conn := tcpserver.InitiateConnection()
	defer conn.Close()

	for {

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

		returnBytes, _ := tcpserver.EncodeLogSLTP(testLog)

		tcpserver.TransmitMessage(conn, &returnBytes)

	}

}
