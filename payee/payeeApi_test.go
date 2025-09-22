package payee

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/payees", PayeePostAPI)
	return mux
}

func TestPayeePostAPISuccess(t *testing.T) {
	mux := setupMux()

	payload := map[string]interface{}{
		"name":           "Abdc",
		"code":           "1263",
		"account_number": 1234767891,
		"ifsc":           "CBIN0123456",
		"bank":           "CBI",
		"email":          "abdc@example.com",
		"mobile":         9876543290,
		"category":       "Employee",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestPayeePostAPIInvalidJSON(t *testing.T) {
	mux := setupMux()

	req := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBufferString("{bad json}"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}
