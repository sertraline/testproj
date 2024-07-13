package validators

import (
	"errors"
	"net/http"

	models "github.com/sertraline/testproj/database/models"
)

// chi router сериализует JSON данные в запросе в эту модель
type OrderCreateRequest struct {
	ExchangeName string              `json:"exchange_name"`
	Pair         string              `json:"pair"`
	Asks         []models.DepthOrder `json:"asks"`
	Bids         []models.DepthOrder `json:"bids"`
}

// валидация входящего запроса
func (u *OrderCreateRequest) Bind(r *http.Request) error {
	if u.ExchangeName == "" {
		return errors.New("exchange name is missing")
	}

	if u.Pair == "" {
		return errors.New("pair is missing")
	}

	if len(u.Asks) == 0 {
		return errors.New("asks is missing. Expected: asks: [{price: 1, base_qty: 1}]}")
	}

	if len(u.Bids) == 0 {
		return errors.New("bids is missing. Expected: bids: [{price: 1, base_qty: 1}]}")
	}

	return nil
}
