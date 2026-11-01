package database

import (
	"context"
	"log"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/joho/godotenv"
)

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
		log.Fatal("Clickhouse Error: ", err)
	}

	return conn, err

}

func DbInit(ctx context.Context, conn driver.Conn) {

	err := conn.Exec(ctx, "CREATE TABLE IF NOT EXISTS LOGS (LogID UUID PRIMARY KEY, DeviceID UUID, SourceAddress IPv4, ProductName String, ServiceName String, Message String, MessageHash String, ScrapeTimestamp DateTime, RecvTimeStamp DateTime);")
	if err != nil {
		log.Fatal(err)
	}

}
