package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	errResp "github.com/sertraline/testproj/errors"
	services "github.com/sertraline/testproj/services"
	validators "github.com/sertraline/testproj/validators"
)

func SaveOrderHistory(w http.ResponseWriter, r *http.Request) {
	// URLParam получает именованные URL запросы
	client_name := chi.URLParam(r, "client_name")
	if client_name == "" {
		render.Render(w, r, errResp.ErrNotFound)
		return
	}

	// валидация и сериализация запроса (используется в chi фреймворке)
	data := &validators.OrderHistoryCreateRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, errResp.ErrInvalidRequest(err))
		return
	}

	// данные приходят в двух отдельных объектах {order{}, client{}}
	data.Client.ClientName = client_name

	// бизнес-логика
	userData, err := services.SaveOrderHistory(data)
	if err != nil {
		render.Render(w, r, errResp.ErrInvalidRequest(err))
	}

	// chi render сериализует данные в JSON
	render.JSON(w, r, userData)
}

func GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	client_name := chi.URLParam(r, "client_name")
	if client_name == "" {
		render.Render(w, r, errResp.ErrNotFound)
		return
	}

	ename := chi.URLParam(r, "exchange_name")
	if ename == "" {
		render.Render(w, r, errResp.ErrNotFound)
		return
	}
	// Query получает URL параметры
	pair := r.URL.Query().Get("pair")
	if pair == "" {
		render.Render(w, r, errResp.ErrNotFound)
		return
	}

	label := r.URL.Query().Get("label")
	if label == "" {
		render.Render(w, r, errResp.ErrNotFound)
		return
	}

	userData, err := services.GetOrderHistory(client_name, label, pair, ename)
	if err != nil {
		render.Render(w, r, errResp.ErrInvalidRequest(err))
	}

	render.JSON(w, r, userData)
}
