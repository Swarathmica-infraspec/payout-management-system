package payee

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

var store PayeeRepository

func initStore() PayeeRepository {
	if store != nil {
		return store
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	store = PayeeDB(db)
	return store
}

func cleanDB(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE payees RESTART IDENTITY CASCADE")
	return err
}

func setupMux(t *testing.T) *http.ServeMux {
	store := initStore()

	payeeDb, ok := store.(*payeeDB)
    if !ok {
        t.Fatalf("store is not *payeeDB")
    }

    if err := cleanDB(payeeDb.db); err != nil {
        t.Fatalf("failed to clean DB: %v", err)
    }

	mux := http.NewServeMux()
	mux.HandleFunc("/payees", PayeePostAPI(store))
	mux.HandleFunc("/payees/list", PayeeGetAPI(store))
	mux.HandleFunc("/payees/", PayeeGetOneAPI(store))
	mux.HandleFunc("/payees/update/", PayeeUpdateAPI(store))
	mux.HandleFunc("/payees/delete/", PayeeDeleteAPI(store))

	return mux
}

func TestPayeePostAPISuccess(t *testing.T) {
	mux := setupMux(t)

	payload := map[string]interface{}{
		"name":           "Abdc",
		"code":           "1262",
		"account_number": 1234767893,
		"ifsc":           "CBIN0123456",
		"bank":           "CBI",
		"email":          "abcd@example.com",
		"mobile":         9876543292,
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
	mux := setupMux(t)

	req := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBufferString("{bad json}"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
}

func TestPayeeGetAPISuccess(t *testing.T) {
	mux := setupMux(t)

	req := httptest.NewRequest(http.MethodGet, "/payees/list", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestPayeeUpdateAPI(t *testing.T) {
	mux := setupMux(t)

	payee := map[string]interface{}{
		"name":           "def",
		"code":           "131",
		"account_number": 1234567090,
		"ifsc":           "SBIN0001111",
		"bank":           "SBI",
		"email":          "def@example.com",
		"mobile":         9876513210,
		"category":       "Employee",
	}
	body, _ := json.Marshal(payee)
	req := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("failed to create payee, got status %d", w.Code)
	}

	updatePayee := map[string]interface{}{
		"name":           "ghhi",
		"code":           "131",
		"account_number": 1234567990,
		"ifsc":           "SBIN0001111",
		"bank":           "SBI",
		"email":          "ghhi@example.com",
		"mobile":         9806517210,
		"category":       "Employee",
	}
	updateBody, _ := json.Marshal(updatePayee)
	req2 := httptest.NewRequest(http.MethodPut, "/payees/update/1", bytes.NewBuffer(updateBody))
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w2.Code, w2.Body.String())
	}
}
func TestPayeeDeleteAPI(t *testing.T) {
	mux := setupMux(t)

	payee := map[string]interface{}{
		"name":           "adef",
		"code":           "1211",
		"account_number": 1134567090,
		"ifsc":           "SBIN0001111",
		"bank":           "SBI",
		"email":          "adef@example.com",
		"mobile":         9876503210,
		"category":       "Employee",
	}
	body, _ := json.Marshal(payee)

	req := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("failed to create payee, got status %d, body=%s", w.Code, w.Body.String())
	}

	req2 := httptest.NewRequest(http.MethodDelete, "/payees/delete/1", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w2.Code, w2.Body.String())
	}

	req3 := httptest.NewRequest(http.MethodGet, "/payees/1", nil)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, req3)

	if w3.Code != http.StatusNotFound {
		t.Fatalf("expected status %d after delete, got %d, body=%s", http.StatusNotFound, w3.Code, w3.Body.String())
	}
}
