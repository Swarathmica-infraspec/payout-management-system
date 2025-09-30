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
func TestPayeePostAPIDuplicate(t *testing.T) {
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

	req1 := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBuffer(body))
	w1 := httptest.NewRecorder()
	mux.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)
	req2 := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBuffer(body))
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusConflict, w2.Code)

	expected := `{"error":"DB insertion failed"}`

	assert.JSONEq(t, expected, w2.Body.String())

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
