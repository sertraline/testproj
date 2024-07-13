package validators

import (
	"errors"
	"net/http"
	"time"

	"github.com/sertraline/testproj/database/models"
)

// входящие JSON данные будут сериализованы в эту модель
type OrderHistoryCreateRequest struct {
	*models.HistoryOrder `json:"order"`
	*models.Client       `json:"client"`
}

// валидация входящего запроса
func (u *OrderHistoryCreateRequest) Bind(r *http.Request) error {
	if (&models.HistoryOrder{} == u.HistoryOrder || u.HistoryOrder == nil) {
		return errors.New("history_order is missing. Expected: order{}")
	}
	if (&models.Client{} == u.Client || u.Client == nil) {
		return errors.New("client is missing. Expected: client{}")
	}

	if u.HistoryOrder.Side == "" {
		return errors.New("side is missing")
	}

	if u.HistoryOrder.Type == "" {
		return errors.New("type is missing")
	}

	if u.HistoryOrder.AlgorithmNamePlaced == "" {
		return errors.New("algorithm_name_placed is missing")
	}

	if time.Time.IsZero(u.HistoryOrder.TimePlaced) {
		return errors.New("time cannot be zero")
	}

	if u.Client.ExchangeName == "" {
		return errors.New("exchange_name is missing")
	}

	if u.Client.Label == "" {
		return errors.New("label is missing")
	}

	if u.Client.Pair == "" {
		return errors.New("pair is missing")
	}

	return nil
}
