package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "github.com/sertraline/testproj/controllers"
)

func TestGetOrderHistoryHandler(t *testing.T) {

	// это null-terminated строка
	expected := `{"status":"Resource not found."}
`

	req := httptest.NewRequest(http.MethodGet, "/clients/John+Doe/ig?label=test&pair=USD/JPY", nil)

	w := httptest.NewRecorder()

	controllers.GetOrderHistory(w, req)

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

func TestPostOrderHistoryHandler(t *testing.T) {

	// это null-terminated строка
	expected := `{"status":"Resource not found."}
`

	req := httptest.NewRequest(http.MethodPost, "/clients/John+Doe", nil)

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
