package database

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type LogSchema struct {
	LogID           uuid.UUID
	DeviceID        uuid.UUID
	SourceAddress   net.IP
	ProductName     []byte
	ServiceName     []byte
	Message         []byte
	MessageHash     [32]byte
	ScrapeTimestamp time.Time
	RecvTimestamp   time.Time
}

func DbConnect(ctx context.Context) (driver.Conn, error) {

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{os.Getenv("CLICKHOUSE_HOSTNAME")},
		Auth: clickhouse.Auth{
			Database: "SIEMSeasy",
			Username: os.Getenv("CLICKHOUSE_USERNAME"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.Ping(ctx); err != nil {
		log.Fatal("ClickHouse Error: ", err)
	}

	return conn, err
}

func DbInit(ctx context.Context, conn driver.Conn) {

	err := conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS LOGS (LogID UUID, DeviceID UUID, SourceAddress IPv4, ProductName String, ServiceName String, Message String, MessageHash FixedString(32), ScrapeTimestamp DateTime, RecvTimestamp DateTime) ENGINE = ReplacingMergeTree ORDER BY (ScrapeTimestamp, ServiceName, MessageHash);")
	if err != nil {
		log.Fatal(err)
	}

}

func DbInsertLog(ctx context.Context, conn driver.Conn, logToInsert []LogSchema) {

	batch, err := conn.PrepareBatch(ctx, "INSERT INTO LOGS")
	if err != nil {
		log.Fatal("ClickHouse Error:", err)
	}
	defer batch.Close()

	for _, singleLog := range logToInsert {
		err := batch.AppendStruct(&singleLog)
		if err != nil {
			log.Fatal("ClickHouse Error:", err)
		}
	}

	err = batch.Send()
	if err != nil {
		log.Fatal("ClickHouse Error:", err)
	}

}
