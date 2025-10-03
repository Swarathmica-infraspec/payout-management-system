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

func TestPayeeGetAPI(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/payees/list", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []PayeeGETResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)

	got := resp[0]
	assert.Equal(t, payload["name"], got.BeneficiaryName)
	assert.Equal(t, payload["code"], got.BeneficiaryCode)
	assert.Equal(t, payload["account_number"], got.AccNo)
	assert.Equal(t, payload["ifsc"], got.IFSC)
	assert.Equal(t, payload["bank"], got.BankName)
	assert.Equal(t, payload["email"], got.Email)
	assert.Equal(t, payload["mobile"], got.Mobile)
	assert.Equal(t, payload["category"], got.PayeeCategory)

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

	req3 := httptest.NewRequest(http.MethodGet, "/payees/"+strconv.Itoa(inserted.ID), nil)
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

	assert.Equal(t, http.StatusCreated, w.Code, "POST to create payee should succeed")

	type PostResponse struct {
		ID int `json:"id"`
	}
	var inserted PostResponse
	err := json.Unmarshal(w.Body.Bytes(), &inserted)
	assert.NoError(t, err, "failed to unmarshal create response")

	url := "/payees/delete/" + strconv.Itoa(inserted.ID)
	req2 := httptest.NewRequest(http.MethodDelete, url, nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code, "DELETE should return 200 OK")

	expectedDeleteResp := `{"status":"deleted"}`
	assert.JSONEq(t, expectedDeleteResp, w2.Body.String(), "DELETE response should match JSON")

	urlGet := "/payees/" + strconv.Itoa(inserted.ID)
	req3 := httptest.NewRequest(http.MethodGet, urlGet, nil)
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusNotFound, w3.Code, "GET after delete should return 404 Not Found")
}
