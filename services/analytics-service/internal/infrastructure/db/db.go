package db

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Config struct {
	host     string
	port     string
	name     string
	user     string
	password string
}

func Connect(ctx context.Context, host, port, name, user, password string) (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{
			host + ":" + port,
		},
		Auth: clickhouse.Auth{
			Database: name,
			Username: user,
			Password: password,
		},
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
