package expense

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ExpenseRequest struct {
	Title        string `json:"title"`
	Amount       int    `json:"amount"`
	DateIncurred string `json:"dateIncurred"`
	Category     string `json:"category"`
	Notes        string `json:"notes"`
	PayeeID      int    `json:"payeeID"`
	ReceiptURI   string `json:"receiptURI"`
}

type ExpenseGETResponse struct {
	Title        string `json:"title"`
	Amount       int    `json:"amount"`
	DateIncurred string `json:"dateIncurred"`
	Category     string `json:"category"`
	Notes        string `json:"notes"`
	PayeeID      int    `json:"payeeID"`
	ReceiptURI   string `json:"receiptURI"`
}

func ExpensePostAPI(store *ExpensePostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data ExpenseRequest
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		e, err := NewExpense(data.Title, float64(data.Amount), data.DateIncurred, data.Category, data.Notes, data.PayeeID, data.ReceiptURI)
		if err != nil {
			http.Error(w, "Invalid payee data", http.StatusBadRequest)
			return
		}

		id, err := store.Insert(context.Background(), e)
		if err != nil {
			http.Error(w, "DB insertion failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
	}
}

func ExpenseGetAPI(store *ExpensePostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		expenses, err := store.List(context.Background())
		if err != nil {
			http.Error(w, "DB query failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var resp []ExpenseGETResponse
		for _, e := range expenses {
			resp = append(resp, ExpenseGETResponse{
				Title:        e.title,
				Amount:       int(e.amount),
				DateIncurred: e.dateIncurred,
				Category:     e.category,
				Notes:        e.notes,
				PayeeID:      e.payeeID,
				ReceiptURI:   e.receiptURI,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func ExpenseGetOneAPI(store *ExpensePostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		e, err := store.GetByID(context.Background(), id)
		if err != nil {
			http.Error(w, "record not found", http.StatusNotFound)
			return
		}

		resp := ExpenseGETResponse{
			Title:        e.title,
			Amount:       int(e.amount),
			DateIncurred: e.dateIncurred,
			Category:     e.category,
			Notes:        e.notes,
			PayeeID:      e.payeeID,
			ReceiptURI:   e.receiptURI,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
