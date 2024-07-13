package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	models "github.com/sertraline/testproj/database/models"
	"github.com/sertraline/testproj/validators"
	"github.com/stretchr/testify/require"
)

func TestSaveOrderHistory(t *testing.T) {
	conn, err := GetNativeConnection(nil, nil, &clickhouse.Compression{
		Method: clickhouse.CompressionLZ4,
	})
	ctx := context.Background()
	require.NoError(t, err)
	fmt.Println("Connection established")

	const ddl = `
		CREATE TABLE OrderHistory (
			client_name				String,
			exchange_name			String,
			label					String,
			pair					String,
			side					String,
			type					String,
			base_qty				Float64,
			price					Float64,
			algorithm_name_placed	String,
			lowest_sell_prc			Float64,
			highest_buy_prc			Float64,
			commission_quote_qty	Float64,
			time_placed				DateTime
		) ENGINE = MergeTree ORDER BY (client_name);
		`
	defer func() {
		fmt.Println("Drop tables")
		conn.Exec(ctx, "DROP TABLE IF EXISTS OrderHistory")
	}()

	fmt.Println("Create OrderHistory")
	require.NoError(t, conn.Exec(ctx, ddl))

	fmt.Println("Initialize data struct")
	data := validators.OrderHistoryCreateRequest{
		HistoryOrder: &models.HistoryOrder{
			ClientName:          "test",
			ExchangeName:        "test",
			Label:               "test",
			Pair:                "test",
			Side:                "test",
			Type:                "test",
			BaseQty:             1,
			Price:               1,
			AlgorithmNamePlaced: "test",
			LowestSellPrc:       1,
			HighestBuyPrc:       1,
			CommissionQuoteQty:  1,
			TimePlaced:          time.Now(),
		},
		Client: &models.Client{
			ClientName:   "test",
			ExchangeName: "test",
			Label:        "test",
			Pair:         "test",
		},
	}

	query := `
		INSERT INTO OrderHistory 
			(
			  client_name, exchange_name, label, pair,
			  side, type, base_qty, price, algorithm_name_placed,
			  lowest_sell_prc, highest_buy_prc,
			  commission_quote_qty, time_placed
			) 
			VALUES (
			  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
			);
		`

	ho := data.HistoryOrder
	cl := data.Client

	fmt.Println("Pack params")
	params := []any{
		cl.ClientName, cl.ExchangeName, cl.Label, cl.Pair,
		ho.Side, ho.Type, ho.BaseQty, ho.Price, ho.AlgorithmNamePlaced,
		ho.LowestSellPrc, ho.HighestBuyPrc, ho.CommissionQuoteQty,
		ho.TimePlaced,
	}

	fmt.Println("Insert record")
	require.NoError(t, conn.Exec(ctx, query, params...))
}

func TestGetOrderHistory(t *testing.T) {
	conn, err := GetNativeConnection(nil, nil, &clickhouse.Compression{
		Method: clickhouse.CompressionLZ4,
	})
	ctx := context.Background()
	require.NoError(t, err)
	fmt.Println("Connection established")

	const ddl = `
		CREATE TABLE OrderHistory (
			client_name				String,
			exchange_name			String,
			label					String,
			pair					String,
			side					String,
			type					String,
			base_qty				Float64,
			price					Float64,
			algorithm_name_placed	String,
			lowest_sell_prc			Float64,
			highest_buy_prc			Float64,
			commission_quote_qty	Float64,
			time_placed				DateTime
		) ENGINE = MergeTree ORDER BY (client_name);
		`
	defer func() {
		fmt.Println("Drop tables")
		conn.Exec(ctx, "DROP TABLE IF EXISTS OrderHistory")
	}()

	require.NoError(t, conn.Exec(ctx, ddl))

	ins := `
		INSERT INTO OrderHistory (
			client_name, exchange_name, label, pair, side, type, base_qty, price,
			algorithm_name_placed, lowest_sell_prc, highest_buy_prc, commission_quote_qty, time_placed
		) VALUES ('test', 'test', 'test', 'test', 'test', 'test', 500, 0.0063, 'test', 0.0062, 0.0065, 0.0035, now64());
	`
	fmt.Println("Insert test record")
	require.NoError(t, conn.Exec(ctx, ins))

	fmt.Println("Initialize data")
	params := []any{"test", "test", "test", "test"}

	query := `
		SELECT * FROM OrderHistory WHERE 
		multiSearchAnyCaseInsensitiveUTF8(client_name, [$1]) 
		AND multiSearchAnyCaseInsensitiveUTF8(exchange_name, [$2])
		AND multiSearchAnyCaseInsensitiveUTF8(label, [$3])
		AND multiSearchAnyCaseInsensitiveUTF8(pair, [$4]) SETTINGS use_query_cache = true;
	`

	fmt.Println("Fetch rows")
	rows, err := conn.Query(context.Background(), query, params...)
	require.NoError(t, err)

	fmt.Println("Fill result with rows")
	orders := make([]*models.HistoryOrder, 0)
	for rows.Next() {
		var ho = &models.HistoryOrder{}
		require.NoError(t, rows.ScanStruct(ho))

		orders = append(orders, ho)
		fmt.Println(ho)
	}
	rows.Close()

	fmt.Printf("Result struct: %#v\n", orders)
	require.NotEmpty(t, orders)
}
