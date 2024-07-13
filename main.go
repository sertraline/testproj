package main

import (
	"fmt"
	"net/http"

	db "github.com/sertraline/testproj/database"
	//models "github.com/sertraline/testproj/database/models"
	controllers "github.com/sertraline/testproj/controllers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()

	// стандартные мидлвари для chi роутера
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // максимально допустимое значение
	}))

	// создание и инициализация таблиц
	db.InitDatabase()

	r.Group(func(r chi.Router) {
		r.Route("/order_book", func(r chi.Router) {
			r.Get("/{exchange_name}", controllers.GetOrderBook)
			r.Post("/", controllers.SaveOrderBook)
		})
		r.Route("/clients", func(r chi.Router) {
			r.Get("/{client_name}/{exchange_name}", controllers.GetOrderHistory)
			r.Post("/{client_name}", controllers.SaveOrderHistory)
		})
	})

	// conn, err := db.GetConn()
	// if err != nil {
	// 	panic("Failed to initialize database connection")
	// }

	// orders := models.OrderBook{}
	// if err := conn.QueryRow(context.Background(), "SELECT * FROM OrderBook").ScanStruct(&orders); err != nil {
	// 	fmt.Println(err)
	// 	panic("Failed to scan, %e")
	// }
	// fmt.Println("result", orders)

	// his := models.HistoryOrder{}
	// if err := conn.QueryRow(context.Background(), "SELECT * FROM OrderHistory").ScanStruct(&his); err != nil {
	// 	fmt.Println(err)
	// 	panic("Failed to scan, %e")
	// }
	// fmt.Println("result", his)

	// client := models.Client{
	// 	ClientName:   his.ClientName,
	// 	ExchangeName: his.ExchangeName,
	// 	Label:        his.Label,
	// 	Pair:         his.Pair,
	// }
	// if err := conn.QueryRow(context.Background(), "SELECT client_name, exchange_name, label, pair FROM OrderHistory").ScanStruct(&client); err != nil {
	// 	fmt.Println(err)
	// 	panic("Failed to scan, %e")
	// }
	// fmt.Println("result", client)
	addr := ":3333"
	fmt.Printf("Starting server on %v\n", addr)
	http.ListenAndServe(addr, r)
}
