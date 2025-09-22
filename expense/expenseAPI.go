package expense

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

func ExpensePostAPI(w http.ResponseWriter, r *http.Request) {

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

	store := NewPostgresExpenseDB(db)


	type req struct {
		Title        string `json:"title"`
		Amount       int    `json:"amount"`
		DateIncurred string `json:"dateIncurred"`
		Category     string `json:"category"`
		Notes        string `json:"notes"`
		PayeeID      int    `json:"payeeID"`
		ReceiptURI   string `json:"receiptURI"`
	}

	var data req

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	e, err := NewExpense(data.Title, float64(data.Amount), data.DateIncurred, data.Category, data.Notes, data.PayeeID, data.ReceiptURI)
	if err != nil {
		fmt.Println("Structure creation failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = store.Insert(context.Background(), e)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Insertion failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received POST request with message")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"id": 1,
	}); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
