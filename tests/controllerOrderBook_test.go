package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "github.com/sertraline/testproj/controllers"
)

func TestGetOrderBookHandler(t *testing.T) {

	// это null-terminated строка
	expected := `{"status":"Resource not found."}
`

	req := httptest.NewRequest(http.MethodGet, "/order_book/test?pair=USD/JPY", nil)

	w := httptest.NewRecorder()

	controllers.GetOrderBook(w, req)

	res := w.Result()

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {

		t.Errorf("Error: %v", err)

	}

	if string(data) != expected {
		t.Errorf("Expected %s found but got [%v]", string(expected), string(data))
	}

}

func TestPostOrderBookHandler(t *testing.T) {

	// это null-terminated строка
	expected := `{"status":"Resource not found."}
`

	req := httptest.NewRequest(http.MethodPost, "/order_book/", nil)

	w := httptest.NewRecorder()

	controllers.GetOrderBook(w, req)

	res := w.Result()

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	if err != nil {

		t.Errorf("Error: %v", err)

	}

	if string(data) != expected {
		t.Errorf("Expected %s found but got [%v]", string(expected), string(data))
	}

}
