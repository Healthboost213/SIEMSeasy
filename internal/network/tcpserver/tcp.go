package tcpserver

import (
	"fmt"
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

	for {
		bytesRead, err := conn.Read(readBuffer)
		if err != nil {
			log.Println("Network Error:", err)
		}
		if bytesRead == 0 {
			break
		}
		fmt.Println("Bytes Read:", bytesRead)
		fmt.Println("Message:", readBuffer[:bytesRead])
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

func TransmitLog(conn *net.TCPConn, sendData *[]byte) {

	n, err := conn.Write(*sendData)
	if err != nil {
		log.Println("Network Error:", err)
	}
	fmt.Println("Bytes Written:", n)

}
