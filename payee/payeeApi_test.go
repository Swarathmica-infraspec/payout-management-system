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
			wantJSON: `{"error":"Payee with the same beneficiary code already exists"}`,
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
			wantJSON: `{"error":"Payee with the same account number already exists"}`,
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
			wantJSON: `{"error":"Payee with the same email already exists"}`,
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
			wantJSON: `{"error":"Payee with the same mobile already exists"}`,
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

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
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

		{"sort by name ASC", "?sort_by=name&sort_order=ASC", []string{"Abdc", "Alice", "Bob", "Charlie"}, 4},
		{"sort by name DESC", "?sort_by=name&sort_order=DESC", []string{"Charlie", "Bob", "Alice", "Abdc"}, 4},
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
			// order has to be checked for sort-related tests
			sortTests := []string{
				"sort by name ASC",
				"sort by name DESC",
				"default sort (id ASC)",
			}

			if contains(sortTests, tt.name) {
				assert.Equal(t, tt.wantNames, gotNames)
			} else {
				assert.ElementsMatch(t, tt.wantNames, gotNames)
			}
		})
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

	assert.Equal(t, http.StatusCreated, w.Code, "first POST should succeed")

	type PostResponse struct {
		ID int `json:"id"`
	}
	var inserted PostResponse
	err := json.Unmarshal(w.Body.Bytes(), &inserted)
	assert.NoError(t, err, "failed to unmarshal create response")

	updatePayee := map[string]interface{}{
		"name":           "ghhi",
		"code":           "131",
		"account_number": 1234567990,
		"ifsc":           "SBIN0002222",
		"bank":           "SBI Updated",
		"email":          "ghhi@example.com",
		"mobile":         9806517210,
		"category":       "Vendor",
	}
	updateBody, _ := json.Marshal(updatePayee)
	req2 := httptest.NewRequest(http.MethodPut, "/payees/update/"+strconv.Itoa(inserted.ID), bytes.NewBuffer(updateBody))
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code, "update should return 200 OK")

	expectedUpdateResp := `{"status":"updated"}`
	assert.JSONEq(t, expectedUpdateResp, w2.Body.String(), "update response should match JSON")

	req3 := httptest.NewRequest(http.MethodGet, "/payees/list?code=131", nil)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code, "GET after update should return 200 OK")

	var got PayeeGETResponse
	err = json.Unmarshal(w3.Body.Bytes(), &got)
	assert.NoError(t, err, "failed to unmarshal get response")

	assert.Equal(t, updatePayee["name"], got.BeneficiaryName, "name should match updated value")
	assert.Equal(t, updatePayee["code"], got.BeneficiaryCode, "code should match updated value")
	assert.Equal(t, updatePayee["account_number"], got.AccNo, "account number should match updated value")
	assert.Equal(t, updatePayee["ifsc"], got.IFSC, "IFSC should match updated value")
	assert.Equal(t, updatePayee["bank"], got.BankName, "bank should match updated value")
	assert.Equal(t, updatePayee["email"], got.Email, "email should match updated value")
	assert.Equal(t, updatePayee["mobile"], got.Mobile, "mobile should match updated value")
	assert.Equal(t, updatePayee["category"], got.PayeeCategory, "category should match updated value")
}
