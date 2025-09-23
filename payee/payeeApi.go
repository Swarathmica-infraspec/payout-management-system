package payee

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func PayeePostAPI(w http.ResponseWriter, r *http.Request) {

	dsn := os.Getenv("TEST_DATABASE_URL")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	store := PayeeDB(db)

	type req struct {
		Name     string `json:"name"`
		Code     string `json:"code"`
		AccNo    int    `json:"account_number"`
		IFSC     string `json:"ifsc"`
		Bank     string `json:"bank"`
		Email    string `json:"email"`
		Mobile   int    `json:"mobile"`
		Category string `json:"category"`
	}

	var data req

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusBadRequest)
		return
	}

	p, err := NewPayee(data.Name, data.Code, data.AccNo, data.IFSC, data.Bank, data.Email, data.Mobile, data.Category)
	if err != nil {
		fmt.Println("Structure creation failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = store.Insert(context.Background(), p)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Insertion failed")
		w.WriteHeader(http.StatusConflict)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": "Payee cannot be created with duplicate values",
		}); err != nil {
			log.Printf("Failed to encode JSON response: %v", err)
			return
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"id": 1,
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}
