package models

var InitOrderBook = `
	CREATE TABLE IF NOT EXISTS OrderBook (
		id 			Int64,
		exchange 	String,
		pair		String,
		asks		Array(Tuple(Float64, Float64)),
		bids		Array(Tuple(Float64, Float64)),
	) ENGINE = MergeTree ORDER BY (id);
`

type DepthOrder struct {
	Price   float64 `json:"price"`
	BaseQty float64 `json:"base_qty"`
}

// Пример:
// INSERT INTO OrderBook (id, exchange, pair, asks, bids) VALUES (1, 'ig', 'USD/JPY', [tuple(0.0062, 10.000)], [tuple(0.0065, 1000.0)]);
// INSERT INTO OrderBook (id, exchange, pair, asks, bids) VALUES (2, 'ex', 'EUR/USD', [tuple(1.08, 4000.0)], [tuple(1.07, 6000.0)]);
//
//		┌─id─┬─exchange─┬─pair────┬─asks──────────┬─bids──────────┐
//	 1. │  2 │ ex   │ EUR/USD │ [(1.08,4000)] │ [(1.07,6000)] │
//	    └────┴──────────┴─────────┴───────────────┴───────────────┘
//	    ┌─id─┬─exchange─┬─pair────┬─asks──────────┬─bids────────────┐
//	 2. │  1 │ ig   │ USD/JPY │ [(0.0062,10)] │ [(0.0065,1000)] │
//	    └────┴──────────┴─────────┴───────────────┴─────────────────┘
type OrderBook struct {
	Id       int64        `ch:"id"`
	Exchange string       `ch:"exchange"`
	Pair     string       `ch:"pair"`
	Asks     []DepthOrder `ch:"asks"`
	Bids     []DepthOrder `ch:"bids"`
}
