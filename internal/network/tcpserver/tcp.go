package tcpserver

import (
	"fmt"
	"io"
	"log"
	"net"
)

// Server Listener

func ListenIncoming() {

	tcpAddress, err := net.ResolveTCPAddr("tcp4", ":6205")
	if err != nil {
		log.Fatal("Network Error:", err)
	}

	netListener, err := net.ListenTCP("tcp4", tcpAddress)
	if err != nil {
		log.Println("Network Error:", err)
	}
	defer netListener.Close()

	for {

		currentConnection, err := netListener.Accept()
		if err != nil {
			log.Println("Network Error:", err)
		}
		go requestHandler(currentConnection)

	}

}

func requestHandler(conn net.Conn) {

	readBuffer := make([]byte, 4096)
	var partial []byte

	for {

		bytesRead, err := conn.Read(readBuffer)

		// Error Conditions
		if err != nil && err != io.EOF {
			log.Println("Network Error:", err)
			break
		}
		if err == io.EOF {
			break
		}
		if bytesRead == 0 {
			break
		}

		// Message Parsing
		message, partial := parseTCPStream(readBuffer[:bytesRead], partial)
		fmt.Printf("Message: %s\n", message)
		fmt.Printf("Partial: %s\n", partial)

	}

}

// Worker Connection

func InitiateConnection() *net.TCPConn {

	tcpAddress, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:6205")
	if err != nil {
		log.Println("Network Error:", err)
	}

	netConnector, err := net.DialTCP("tcp4", nil, tcpAddress)
	if err != nil {
		log.Println("Network Error:", err)
	}

	return netConnector

}

func TransmitMessage(conn *net.TCPConn, sendData *[]byte) {

	_, err := conn.Write(*sendData)
	if err != nil {
		log.Println("Network Error:", err)
	}

}
