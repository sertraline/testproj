package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	db "github.com/sertraline/testproj/database"
	models "github.com/sertraline/testproj/database/models"
	validators "github.com/sertraline/testproj/validators"
)

func SaveOrderBook(data *validators.OrderCreateRequest) (int, error) {
	// Я не нашел autoincrement в ClickHouse, поэтому я определяю id записи через Count.
	//

	query := `
	INSERT INTO OrderBook (id, exchange, pair, asks, bids) SELECT
		COUNT(),
		$1,
		$2,
	`

	params := []any{data.ExchangeName, data.Pair}

	// Динамическая распаковка массива с Asks и Bids
	// [
	//   {"price": 1.2, "base_qty": 200},
	//   {"price": 1.3, "base_qty": 800}
	// ] будет представлено как [tuple($3, $4), tuple($5, $6)] в конечном запросе.
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

	err := db.AsyncInsert(query, params...)
	if err != nil {
		fmt.Println(err)
		return 1, err
	}

	return 0, nil
}

func GetOrderBook(exchangeName string, pair string) (models.OrderBook, error) {
	query := `
		SELECT * FROM OrderBook WHERE 
		multiSearchAnyCaseInsensitiveUTF8(exchange, [$1]) 
		AND multiSearchAnyCaseInsensitiveUTF8(pair, [$2]);
	`

	conn, err := db.GetConn()
	if err != nil {
		panic("Failed to initialize database connection")
	}

	orders := models.OrderBook{}
	if err := conn.QueryRow(context.Background(), query, exchangeName, pair).ScanStruct(&orders); err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			orders.Id = -1
			return orders, err
		} else {
			return orders, err
		}
	}
	return orders, nil
}
