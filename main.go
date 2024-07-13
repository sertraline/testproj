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

	// присваивает уникальный идентификатор каждому запросу
	r.Use(middleware.RequestID)
	// логгирует каждый запрос
	r.Use(middleware.Logger)
	// позволяет дочерним тредам паниковать и восстанавливать свою работу
	r.Use(middleware.Recoverer)
	// позволяет использовать именованные URL параметры типа {name}
	r.Use(middleware.URLFormat)

	r.Use(render.SetContentType(render.ContentTypeJSON))
	// CORS настроен на любые источники
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

	addr := ":3333"
	fmt.Printf("Starting server on %v\n", addr)
	http.ListenAndServe(addr, r)
}
