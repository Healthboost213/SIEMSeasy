package tcp_server

import (
	"fmt"
	"log"
	"net"
	"time"
)

// Server Listener

func ListenIncoming() {

	tcpAddress, err := net.ResolveTCPAddr("tcp4", ":6205")
	if err != nil {
		log.Fatal("Network Error:", err)
	}

	netListener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		log.Fatal("Network Error:", err)
	}
	defer netListener.Close()

	for {

		currentConnection, err := netListener.Accept()
		if err != nil {
			log.Fatal("Network Error:", err)
		}
		defer currentConnection.Close()

		go requestHandler(currentConnection)
	}

}

func requestHandler(conn net.Conn) {

	readBuffer := make([]byte, 4096)
	writeBuffer := []byte("PONG\n")

	bytesRead, err := conn.Read(readBuffer)
	if err != nil {
		log.Fatal("Network Error:", err)
	}

	fmt.Println("Bytes Read:", bytesRead)
	fmt.Printf("Message: %s", readBuffer[:bytesRead])

	// simple pong response ever 5 sec | add functionality later
	for {
		_, err = conn.Write(writeBuffer)
		if err != nil {
			log.Printf("Network Error: %s", err)
			break
		}
		fmt.Println("sent PONG")
		dur := time.Duration(5) * time.Second
		time.Sleep(dur)
	}

}
