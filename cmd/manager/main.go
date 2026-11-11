package main

import (
	"SIEMSEasy/internal/database"
	"SIEMSEasy/internal/httpserver"
	"SIEMSEasy/internal/realtime"
	"context"
	"log"
)

func main() {

	ctx := context.Background()

	mainBuffer := realtime.CreateRingBuffer()

	conn, err := database.DbConnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	database.DbInit(ctx, conn)

	httpserver.StartManagerHTTP()

}
