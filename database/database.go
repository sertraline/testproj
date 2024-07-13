package database

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	models "github.com/sertraline/testproj/database/models"
)

func GetConn() (clickhouse.Conn, error) {
	// бд "orders" генерируется в fs/volumes/clickhouse/docker-entrypoint-inidb.d
	addr := "127.0.0.1:9000"
	log.Printf("Connect %s", addr)
	dialCount := 0
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: "orders",
			Username: "default",
			// без пароля
			Password: "",
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			dialCount++
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: true,
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "OrderAPI", Version: "0.0.1"},
			},
		},
	})

	if err != nil {
		return nil, err
	}
	// проверка подключения
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, err
}

func AsyncInsert(query string, params ...any) error {
	conn, err := GetConn()
	if err != nil {
		return err
	}
	ctx := context.Background()

	defer func() {
		conn.Close()
	}()

	if err := conn.Exec(ctx, query, params...); err != nil {
		return err
	}

	return nil
}

func InitDatabase() error {
	if err := AsyncInsert(models.InitOrderBook); err != nil {
		return err
	}
	if err := AsyncInsert(models.InitOrderHistory); err != nil {
		return err
	}
	return nil
}