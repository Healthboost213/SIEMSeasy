package main

import (
	"SIEMSEasy/internal/database"
	"SIEMSEasy/internal/realtime"
	"context"
	"log"
)

func main() {

	var queue realtime.RingBuffer
	queue.Head = 0
	queue.Tail = 0

	ctx := context.Background()

	conn, err := database.DbConnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	database.DbInit(ctx, conn)

}
