package database

import (
	"SIEMSEasy/internal/httpserver"
	"context"
	"log"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/joho/godotenv"
)

type LogSchema struct {
	LogID           [16]byte
	DeviceID        [16]byte
	SourceAddress   [4]byte
	ProductName     string
	ServiceName     string
	Message         string
	MessageHash     [32]byte
	ScrapeTimestamp time.Time
	RecvTimestamp   time.Time
}

func ConvertJsonToDb(jsonFormat httpserver.LogJSON) LogSchema {

	var dbData LogSchema

	// Byte Conversions
	copy(dbData.LogID[:], jsonFormat.LogID)
	copy(dbData.DeviceID[:], jsonFormat.DeviceID)
	copy(dbData.SourceAddress[:], jsonFormat.SourceAddress)
	copy(dbData.MessageHash[:], jsonFormat.MessageHash)

	// String Conversions
	dbData.ProductName = jsonFormat.ProductName
	dbData.ServiceName = jsonFormat.ServiceName
	dbData.Message = jsonFormat.Message

	// Time Conversion
	dbData.ScrapeTimestamp, _ = time.Parse(time.RFC3339Nano, jsonFormat.ScrapeTimestamp)
	dbData.RecvTimestamp = time.Now()

	return dbData

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

	err := conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS LOGS (LogID UUID, DeviceID UUID, SourceAddress IPv4, ProductName String, ServiceName String, Message String, MessageHash String, ScrapeTimestamp DateTime, RecvTimeStamp DateTime) ENGINE = ReplacingMergeTree ORDER BY (MessageHash);")
	if err != nil {
		log.Fatal(err)
	}

}
