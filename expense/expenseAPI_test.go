package expense

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/expenses", ExpensePostAPI)
	return mux
}

func TestExpensePostAPISuccess(t *testing.T) {
	mux := setupMux()

	payload := map[string]interface{}{
		"title":          "Food",
		"amount":         100,
		"dateIncurred":   "2025-09-06",
		"category":       "bill",
		"notes":          "dinner",
		"payeeID":        1,
		"receiptURI":     "/food_bill.jpg",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestExpensePostAPIInvalidJSON(t *testing.T) {
	mux := setupMux()

	req, _ := http.NewRequest("POST", "/expensex", bytes.NewBufferString("{bad json}"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
