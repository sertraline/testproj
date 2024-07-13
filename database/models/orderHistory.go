package models

import "time"

var InitOrderHistory = `
	CREATE TABLE IF NOT EXISTS OrderHistory (
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

// Пример:
// INSERT INTO OrderHistory (client_name, exchange_name, label, pair, side, type, base_qty, price, algorithm_name_placed, lowest_sell_prc, highest_buy_prc, commission_quote_qty, time_placed) VALUES ('John Doe', 'ig', 'Buy', 'USD/JPY', 'LHS', 'lever', '500', 0.0063, 'statistics', 0.0062, 0.0065, 0.0035, now64());
//
//		┌─client_name─┬─exchange_name─┬─label─┬─pair────┬─side─┬─type──┬─base_qty─┬──price─┬─algorithm_name_placed─┬─lowest_sell_prc─┬─highest_buy_prc─┬─commission_quote_qty─┬─────────time_placed─┐
//	 1. │ John Doe    │ ig        │ Buy   │ USD/JPY │ LHS  │ lever │      500 │ 0.0063 │ statistics            │          0.0062 │          0.0065 │               0.0035 │ 2024-07-11 15:08:50 │
//	    └─────────────┴───────────────┴───────┴─────────┴──────┴───────┴──────────┴────────┴───────────────────────┴─────────────────┴─────────────────┴──────────────────────┴─────────────────────┘
type HistoryOrder struct {
	ClientName          string    `ch:"client_name" json:"client_name"`
	ExchangeName        string    `ch:"exchange_name" json:"exchange_name"`
	Label               string    `ch:"label"`
	Pair                string    `ch:"pair"`
	Side                string    `ch:"side"`
	Type                string    `ch:"type"`
	BaseQty             float64   `ch:"base_qty" json:"base_qty"`
	Price               float64   `ch:"price"`
	AlgorithmNamePlaced string    `ch:"algorithm_name_placed" json:"algorithm_name_placed"`
	LowestSellPrc       float64   `ch:"lowest_sell_prc" json:"lowest_sell_prc"`
	HighestBuyPrc       float64   `ch:"highest_buy_prc" json:"highest_buy_prc"`
	CommissionQuoteQty  float64   `ch:"commission_quote_qty" json:"commission_quote_qty"`
	TimePlaced          time.Time `ch:"time_placed" json:"time_placed"`
}

type Client struct {
	ClientName   string `ch:"client_name" json:"client_name"`
	ExchangeName string `ch:"exchange_name" json:"exchange_name"`
	Label        string `ch:"label"`
	Pair         string `ch:"pair"`
}
