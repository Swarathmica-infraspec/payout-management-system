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
	if !ok {
		t.Fatalf("store is not *payeeDB")
	}

	if err := cleanDB(payeeDb.db); err != nil {
		t.Fatalf("failed to clean DB: %v", err)
	}

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

	type Response struct {
		ID int `json:"id"`
	}

	var resp Response
	err := json.Unmarshal([]byte(w.Body.Bytes()), &resp)
	if err != nil {
		t.Fatal("Error unmarshaling JSON:", err)
		return
	}

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusCreated, w.Code, w.Body.String())
	}
	expected := `{"id":1}` + "\n"
	if w.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, w.Body.String())
	}
	if resp.ID != 1 {
		t.Fatalf("The response body should be {\"id\":1}")
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

	resp := w.Body.String()
	expected := "Invalid JSON body\n"

	if resp != expected {
		t.Fatalf("expected body %q, got %q", expected, resp)
	}

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
	if w1.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", w1.Code)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/payees", bytes.NewBuffer(body))
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Fatalf("expected 409 Conflict, got %d, body=%s", w2.Code, w2.Body.String())
	}

	expected := "DB insertion failed: pq: duplicate key value violates unique constraint \"payees_beneficiary_code_key\"\n"
	if w2.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, w2.Body.String())
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
	if wCreate.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d, body=%s", wCreate.Code, wCreate.Body.String())
	}

	req := httptest.NewRequest(http.MethodGet, "/payees/list", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	var resp []PayeeGETResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if len(resp) != 1 {
		t.Fatalf("expected 1 payee, got %d", len(resp))
	}

	got := resp[0]
	if got.BeneficiaryName != payload["name"] {
		t.Errorf("expected name %q, got %q", payload["name"], got.BeneficiaryName)
	}
	if got.BeneficiaryCode != payload["code"] {
		t.Errorf("expected code %q, got %q", payload["code"], got.BeneficiaryCode)
	}
	if got.AccNo != payload["account_number"] {
		t.Errorf("expected account_number %v, got %v", payload["account_number"], got.AccNo)
	}
	if got.IFSC != payload["ifsc"] {
		t.Errorf("expected IFSC %q, got %q", payload["ifsc"], got.IFSC)
	}
	if got.BankName != payload["bank"] {
		t.Errorf("expected bank %q, got %q", payload["bank"], got.BankName)
	}
	if got.Email != payload["email"] {
		t.Errorf("expected email %q, got %q", payload["email"], got.Email)
	}
	if got.Mobile != payload["mobile"] {
		t.Errorf("expected mobile %v, got %v", payload["mobile"], got.Mobile)
	}
	if got.PayeeCategory != payload["category"] {
		t.Errorf("expected category %q, got %q", payload["category"], got.PayeeCategory)
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

	if wCreate.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d, body=%s", wCreate.Code, wCreate.Body.String())
	}

	type CreateResp struct {
		ID int `json:"id"`
	}
	var createResp CreateResp
	if err := json.Unmarshal(wCreate.Body.Bytes(), &createResp); err != nil {
		t.Fatal("Failed to unmarshal create response:", err)
	}

	url := "/payees/" + strconv.Itoa(createResp.ID)
	reqGet := httptest.NewRequest(http.MethodGet, url, nil)
	wGet := httptest.NewRecorder()
	mux.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d, body=%s", wGet.Code, wGet.Body.String())
	}

	var getResp PayeeGETResponse

	if err := json.Unmarshal(wGet.Body.Bytes(), &getResp); err != nil {
		t.Fatal("Failed to unmarshal get response:", err)
	}

	if getResp.ID != createResp.ID {
		t.Fatalf("expected ID %d, got %d", createResp.ID, getResp.ID)
	}
	if getResp.BeneficiaryName != payload["name"] {
		t.Fatalf("expected name %q, got %q", payload["name"], getResp.BeneficiaryName)
	}
	if getResp.BeneficiaryCode != payload["code"] {
		t.Fatalf("expected code %q, got %q", payload["code"], getResp.BeneficiaryCode)
	}
	if getResp.AccNo != payload["account_number"] {
		t.Fatalf("expected accNo %v, got %v", payload["account_number"], getResp.AccNo)
	}
	if getResp.IFSC != payload["ifsc"] {
		t.Fatalf("expected IFSC %q, got %q", payload["ifsc"], getResp.IFSC)
	}
	if getResp.BankName != payload["bank"] {
		t.Fatalf("expected bank %q, got %q", payload["bank"], getResp.BankName)
	}
	if getResp.Email != payload["email"] {
		t.Fatalf("expected email %q, got %q", payload["email"], getResp.Email)
	}
	if getResp.Mobile != payload["mobile"] {
		t.Fatalf("expected mobile %v, got %v", payload["mobile"], getResp.Mobile)
	}
	if getResp.PayeeCategory != payload["category"] {
		t.Fatalf("expected category %q, got %q", payload["category"], getResp.PayeeCategory)
	}
}

func TestPayeeGetOneAPINotFound(t *testing.T) {
	mux := setupMux(t)

	nonExistentID := 9999
	url := "/payees/" + strconv.Itoa(nonExistentID)

	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d, body=%s", w.Code, w.Body.String())
	}

	expected := "record not found\n"
	if w.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, w.Body.String())
	}
}
