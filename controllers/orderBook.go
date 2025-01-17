package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	errResp "github.com/sertraline/testproj/errors"
	services "github.com/sertraline/testproj/services"
	validators "github.com/sertraline/testproj/validators"
)

func SaveOrderBook(w http.ResponseWriter, r *http.Request) {
	// валидация и сериализация запроса (используется в chi фреймворке)
	data := &validators.OrderCreateRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, errResp.ErrInvalidRequest(err))
		return
	}

	// бизнес-логика
	userData, err := services.SaveOrderBook(data)
	if err != nil {
		render.Render(w, r, errResp.ErrInvalidRequest(err))
	}

	// chi render сериализует данные в JSON
	render.JSON(w, r, userData)
}

func GetOrderBook(w http.ResponseWriter, r *http.Request) {
	// URLParam получает именованные URL запросы
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

	userData, err := services.GetOrderBook(ename, pair)
	if err != nil {
		if userData.Id == -1 {
			render.Render(w, r, errResp.ErrNotFound)
		} else {
			render.Render(w, r, errResp.ErrInvalidRequest(err))
		}
	}

	render.JSON(w, r, userData)
}
