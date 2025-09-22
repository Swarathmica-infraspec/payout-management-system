package expense

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var store *ExpensePostgresDB

func initStore() *ExpensePostgresDB {
	if store != nil {
		return store
	}
	dsn := os.Getenv("TEST_DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	store = NewPostgresExpenseDB(db)
	return store
}

func setupMux() *http.ServeMux {
	store := initStore()
	mux := http.NewServeMux()
	mux.HandleFunc("/expenses", ExpensePostAPI(store))
	mux.HandleFunc("/expenses/list", ExpenseGetAPI(store))
	mux.HandleFunc("/expenses/", ExpenseGetOneAPI(store))
	return mux
}

func TestExpensePostAPISuccess(t *testing.T) {
	mux := setupMux()

	payload := map[string]interface{}{
		"title":        "Food",
		"amount":       100,
		"dateIncurred": "2025-09-06",
		"category":     "bill",
		"notes":        "dinner",
		"payeeID":      1,
		"receiptURI":   "/food_bill.jpg",
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

	req, _ := http.NewRequest("POST", "/expenses", bytes.NewBufferString("{bad json}"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestExpenseGetAPISuccess(t *testing.T) {
	mux := setupMux()

	req := httptest.NewRequest(http.MethodGet, "/expenses/list", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
}
