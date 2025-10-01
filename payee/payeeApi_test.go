package payee

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var store PayeeRepository

func initStore() PayeeRepository {
	if store != nil {
		return store
	}
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
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
	require.True(t, ok, "store should be *payeeDB")

	err := cleanDB(payeeDb.db)
	require.NoError(t, err, "failed to clean DB")

	mux := SetupRouter(store)

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

	assert.Equal(t, http.StatusCreated, w.Code)

	expected := `{"id":1}`
	assert.JSONEq(t, expected, w.Body.String())

}

func TestPayeePostAPIInvalidJSON(t *testing.T) {
	mux := setupMux(t)

	req := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBufferString("{bad json}"))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Invalid JSON body"}`
	assert.JSONEq(t, expected, w.Body.String())

}
func TestPayeePostAPIUniqueConstraints(t *testing.T) {
	mux := setupMux(t)

	// Original payee to create first
	original := map[string]interface{}{
		"name":           "Abc",
		"code":           "136",
		"account_number": 1234567890123456,
		"ifsc":           "CBIN0123459",
		"bank":           "CBI",
		"email":          "abc@gmail.com",
		"mobile":         9123456780,
		"category":       "Employee",
	}
	body, _ := json.Marshal(original)

	// Insert original payee
	req1 := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBuffer(body))
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// Table-driven duplicate tests
	tests := []struct {
		name     string
		payload  map[string]interface{}
		wantCode int
		wantJSON string
	}{
		{
			name: "duplicate code",
			payload: map[string]interface{}{
				"name":           "Xyz",
				"code":           "136",
				"account_number": 9876543210123456,
				"ifsc":           "CBIN0123460",
				"bank":           "CBI",
				"email":          "xyz@gmail.com",
				"mobile":         9876543290,
				"category":       "Employee",
			},
			wantCode: http.StatusConflict,
			wantJSON: `{"error":"Payee already exists with the same: code"}`,
		},
		{
			name: "duplicate account",
			payload: map[string]interface{}{
				"name":           "Xyz",
				"code":           "137",
				"account_number": 1234567890123456,
				"ifsc":           "CBIN0123460",
				"bank":           "CBI",
				"email":          "x@gmail.com",
				"mobile":         9876543291,
				"category":       "Employee",
			},
			wantCode: http.StatusConflict,
			wantJSON: `{"error":"Payee already exists with the same: account number"}`,
		},
		{
			name: "duplicate email",
			payload: map[string]interface{}{
				"name":           "Pqr",
				"code":           "138",
				"account_number": 1111111111111111,
				"ifsc":           "CBIN0123461",
				"bank":           "CBI",
				"email":          "abc@gmail.com",
				"mobile":         9876543292,
				"category":       "Employee",
			},
			wantCode: http.StatusConflict,
			wantJSON: `{"error":"Payee already exists with the same: email"}`,
		},
		{
			name: "duplicate mobile",
			payload: map[string]interface{}{
				"name":           "Lmn",
				"code":           "139",
				"account_number": 2222222222222222,
				"ifsc":           "CBIN0123462",
				"bank":           "CBI",
				"email":          "lmn@gmail.com",
				"mobile":         9123456780,
				"category":       "Employee",
			},
			wantCode: http.StatusConflict,
			wantJSON: `{"error":"Payee already exists with the same: mobile"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.JSONEq(t, tt.wantJSON, w.Body.String())
		})
	}
}
