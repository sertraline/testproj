package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/ClickHouse/clickhouse-go/v2"
	models "github.com/sertraline/testproj/database/models"
	"github.com/sertraline/testproj/validators"
	"github.com/stretchr/testify/require"
)

func TestSaveOrderBook(t *testing.T) {
	conn, err := GetNativeConnection(nil, nil, &clickhouse.Compression{
		Method: clickhouse.CompressionLZ4,
	})
	ctx := context.Background()

	require.NoError(t, err)

	fmt.Println("Connection established")
	const ddl = `
		CREATE TABLE OrderBook (
			id 			Int64,
			exchange 	String,
			pair		String,
			asks		Array(Tuple(Float64, Float64)),
			bids		Array(Tuple(Float64, Float64)),
		) Engine MergeTree() ORDER BY tuple()
		`
	defer func() {
		fmt.Println("Drop tables")
		conn.Exec(ctx, "DROP TABLE IF EXISTS OrderBook")
	}()
	fmt.Println("Create OrderBook")
	require.NoError(t, conn.Exec(ctx, ddl))

	fmt.Println("Initialize data struct")
	ask := []models.DepthOrder{{Price: 1, BaseQty: 1}}
	data := validators.OrderCreateRequest{
		ExchangeName: "test",
		Pair:         "test",
		Asks:         ask,
		Bids:         ask,
	}

	query := `
	INSERT INTO OrderBook (id, exchange, pair, asks, bids) SELECT
		COUNT(),
		$1,
		$2,
	`

	fmt.Println("Unpack query")
	params := []any{data.ExchangeName, data.Pair}

	counter := 2
	query = query + "["
	for _, v := range data.Asks {
		counter += 2
		params = append(params, v.Price, v.BaseQty)
		query = query + fmt.Sprintf("tuple($%d, $%d),", counter-1, counter)
	}
	query = strings.Trim(query, ",")
	query = query + "],["

	for _, v := range data.Bids {
		counter += 2
		params = append(params, v.Price, v.BaseQty)
		query = query + fmt.Sprintf("tuple($%d, $%d),", counter-1, counter)
	}
	query = strings.Trim(query, ",")
	query = query + "] "

	query = query + "FROM OrderBook;"

	fmt.Println("Execute query")
	require.NoError(t, conn.Exec(ctx, query, params...))
}

func TestGetOrderBook(t *testing.T) {
	conn, err := GetNativeConnection(nil, nil, &clickhouse.Compression{
		Method: clickhouse.CompressionLZ4,
	})
	ctx := context.Background()
	require.NoError(t, err)

	fmt.Println("Connection established")
	const ddl = `
		CREATE TABLE OrderBook (
			id 			Int64,
			exchange 	String,
			pair		String,
			asks		Array(Tuple(Float64, Float64)),
			bids		Array(Tuple(Float64, Float64)),
		) Engine MergeTree() ORDER BY tuple()
		`
	defer func() {
		fmt.Println("Drop tables")
		conn.Exec(ctx, "DROP TABLE IF EXISTS OrderBook")
	}()
	fmt.Println("Create OrderBook")
	require.NoError(t, conn.Exec(ctx, ddl))

	ins := `
		INSERT INTO OrderBook (id, exchange, pair, asks, bids) VALUES (1, 'test', 'test', [tuple(0.0062, 10.000)], [tuple(0.0065, 1000.0)]);
	`
	fmt.Println("Insert test record")
	require.NoError(t, conn.Exec(ctx, ins))

	query := `
		SELECT * FROM OrderBook WHERE
		multiSearchAnyCaseInsensitiveUTF8(exchange, [$1])
		AND multiSearchAnyCaseInsensitiveUTF8(pair, [$2]) SETTINGS use_query_cache = true;;
	`

	exchangeName := "test"
	pair := "test"

	orders := models.OrderBook{}
	fmt.Println("Fetch row")
	require.NoError(t, conn.QueryRow(context.Background(), query, exchangeName, pair).ScanStruct(&orders))
	fmt.Println("Check result struct")
	require.NotEmpty(t, orders.Exchange)
}
