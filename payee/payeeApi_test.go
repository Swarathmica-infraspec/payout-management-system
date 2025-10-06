package payee

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	db, err := sql.Open("pgx", dsn)
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
	createPayee := func(payload map[string]interface{}) *httptest.ResponseRecorder {
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		return w
	}
	w := createPayee(original)
	assert.Equal(t, http.StatusCreated, w.Code)

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
			wantJSON: `{"error":"Payee already exists with the same: beneficiary code"}`,
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
			w := createPayee(tt.payload)
			assert.Equal(t, tt.wantCode, w.Code)
			assert.JSONEq(t, tt.wantJSON, w.Body.String())
		})
	}
}

func TestPayeeGetOneAPI(t *testing.T) {
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

	reqCreate := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBuffer(body))
	wCreate := httptest.NewRecorder()
	mux.ServeHTTP(wCreate, reqCreate)

	assert.Equal(t, http.StatusCreated, wCreate.Code)

	type CreateResp struct {
		ID int `json:"id"`
	}
	var createResp CreateResp
	err := json.Unmarshal(wCreate.Body.Bytes(), &createResp)
	assert.NoError(t, err)

	url := "/payees/" + strconv.Itoa(createResp.ID)
	reqGet := httptest.NewRequest(http.MethodGet, url, nil)
	wGet := httptest.NewRecorder()
	mux.ServeHTTP(wGet, reqGet)

	assert.Equal(t, http.StatusOK, wGet.Code)

	var getResp PayeeGETResponse
	err = json.Unmarshal(wGet.Body.Bytes(), &getResp)
	assert.NoError(t, err)

	assert.Equal(t, createResp.ID, getResp.ID)
	assert.Equal(t, payload["name"], getResp.BeneficiaryName)
	assert.Equal(t, payload["code"], getResp.BeneficiaryCode)
	assert.Equal(t, payload["account_number"], getResp.AccNo)
	assert.Equal(t, payload["ifsc"], getResp.IFSC)
	assert.Equal(t, payload["bank"], getResp.BankName)
	assert.Equal(t, payload["email"], getResp.Email)
	assert.Equal(t, payload["mobile"], getResp.Mobile)
	assert.Equal(t, payload["category"], getResp.PayeeCategory)
}

func TestPayeeGetOneAPINotFound(t *testing.T) {
	mux := setupMux(t)

	nonExistentID := 9999
	url := "/payees/" + strconv.Itoa(nonExistentID)

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	expected := `{"error":"record not found"}`
	assert.JSONEq(t, expected, w.Body.String())

}

func TestPayeeGetAPI(t *testing.T) {
	mux := setupMux(t)

	payees := []map[string]interface{}{
		{"name": "Alice", "code": "A001", "account_number": 1112345678901324, "ifsc": "HDFC0017890", "bank": "HDFC", "email": "a@example.com", "mobile": 9000000001, "category": "Vendor"},
		{"name": "Bob", "code": "B001", "account_number": 2225678347532479, "ifsc": "SBIN0022345", "bank": "SBI", "email": "b@example.com", "mobile": 9000000002, "category": "Employee"},
		{"name": "Charlie", "code": "C001", "account_number": 3335674839247567, "ifsc": "HDFC0033333", "bank": "HDFC", "email": "c@example.com", "mobile": 9000000003, "category": "Vendor"},
		{"name": "Abdc", "code": "1262", "account_number": 1234767893, "ifsc": "CBIN0123456", "bank": "CBI", "email": "abcd@example.com", "mobile": 9876543292, "category": "Employee"},
	}

	for _, p := range payees {
		body, _ := json.Marshal(p)
		req := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}

	tests := []struct {
		name       string
		query      string
		wantNames  []string
		wantLength int
	}{
		{"list all", "", []string{"Alice", "Bob", "Charlie", "Abdc"}, 4},

		{"filter by bank HDFC", "?bank=HDFC", []string{"Alice", "Charlie"}, 2},
		{"filter by category Employee", "?category=Employee", []string{"Bob", "Abdc"}, 2},
		{"filter by name Alice", "?name=Alice", []string{"Alice"}, 1},
		{"filter by bank HDFC & category Vendor", "?bank=HDFC&category=Vendor", []string{"Alice", "Charlie"}, 2},

		{"sort by name ASC", "?sort_by=beneficiary_name&sort_order=ASC", []string{"Abdc", "Alice", "Bob", "Charlie"}, 4},
		{"sort by name DESC", "?sort_by=beneficiary_name&sort_order=DESC", []string{"Charlie", "Bob", "Alice", "Abdc"}, 4},
		{"default sort (id ASC)", "", []string{"Alice", "Bob", "Charlie", "Abdc"}, 4},

		{"pagination limit 1 offset 0", "?sort_by=id&sort_order=ASC&limit=1&offset=0", []string{"Alice"}, 1},
		{"pagination limit 1 offset 1", "?sort_by=id&sort_order=ASC&limit=1&offset=1", []string{"Bob"}, 1},
		{"pagination limit 2 offset 1", "?sort_by=id&sort_order=ASC&limit=2&offset=1", []string{"Bob", "Charlie"}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/payees/list"+tt.query, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var resp []PayeeGETResponse
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Len(t, resp, tt.wantLength)

			var gotNames []string
			for _, p := range resp {
				gotNames = append(gotNames, p.BeneficiaryName)
			}
			assert.ElementsMatch(t, tt.wantNames, gotNames)
		})
	}
}
