package services

import (
	"context"
	"fmt"
	"strings"

	db "github.com/sertraline/testproj/database"
	models "github.com/sertraline/testproj/database/models"
	validators "github.com/sertraline/testproj/validators"
)

func SaveOrderHistory(data *validators.OrderHistoryCreateRequest) (int, error) {
	// Я не нашел autoincrement в ClickHouse, поэтому я определяю id записи через Count.
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

	fmt.Println(ho, cl)
	params := []any{
		cl.ClientName, cl.ExchangeName, cl.Label, cl.Pair,
		ho.Side, ho.Type, ho.BaseQty, ho.Price, ho.AlgorithmNamePlaced,
		ho.LowestSellPrc, ho.HighestBuyPrc, ho.CommissionQuoteQty,
		ho.TimePlaced,
	}
	err := db.AsyncInsert(query, params...)
	if err != nil {
		fmt.Println(err)
		return 1, err
	}

	return 0, nil
}

func GetOrderHistory(clientName string, label string, pair string, exchangeName string) ([]*models.HistoryOrder, error) {
	// пробелы в имени клиента можно обозначать через +
	clientName = strings.ReplaceAll(clientName, "+", " ")

	query := `
		SELECT * FROM OrderHistory WHERE 
		multiSearchAnyCaseInsensitiveUTF8(client_name, [$1]) 
		AND multiSearchAnyCaseInsensitiveUTF8(exchange_name, [$2])
		AND multiSearchAnyCaseInsensitiveUTF8(label, [$3])
		AND multiSearchAnyCaseInsensitiveUTF8(pair, [$4]) SETTINGS use_query_cache = true;
	`

	conn, err := db.GetConn()
	if err != nil {
		panic("Failed to initialize database connection")
	}

	defer func() {
		conn.Close()
	}()

	rows, err := conn.Query(context.Background(), query, clientName, exchangeName, label, pair)
	if err != nil {
		fmt.Println(err)
		return []*models.HistoryOrder{}, err
	}

	orders := make([]*models.HistoryOrder, 0)
	for rows.Next() {
		var ho = &models.HistoryOrder{}
		if err := rows.ScanStruct(ho); err != nil {
			return []*models.HistoryOrder{}, err
		}

		orders = append(orders, ho)
		fmt.Println(ho)
	}
	rows.Close()

	fmt.Println("result", orders)
	return orders, rows.Err()
}
