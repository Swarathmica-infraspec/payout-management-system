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
	"strconv"

	_ "github.com/lib/pq"
)

func PayeePostAPI(w http.ResponseWriter, r *http.Request) {

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	store := PostgresPayeeDB(db)

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

type PayeeGETResponse struct {
	ID              int    `json:"id"`
	BeneficiaryName string `json:"beneficiary_name"`
	BeneficiaryCode string `json:"beneficiary_code"`
	AccNo           int    `json:"account_number"`
	IFSC            string `json:"ifsc_code"`
	BankName        string `json:"bank_name"`
	Email           string `json:"email"`
	Mobile          int    `json:"mobile"`
	PayeeCategory   string `json:"payee_category"`
}

func PayeeGetAPI(store *PayeePostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payees, err := store.List(context.Background())
		if err != nil {
			http.Error(w, "DB query failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var resp []PayeeGETResponse
		for _, p := range payees {
			resp = append(resp, PayeeGETResponse{
				ID:              p.id,
				BeneficiaryName: p.beneficiaryName,
				BeneficiaryCode: p.beneficiaryCode,
				AccNo:           p.accNo,
				IFSC:            p.ifsc,
				BankName:        p.bankName,
				Email:           p.email,
				Mobile:          p.mobile,
				PayeeCategory:   p.payeeCategory,
			})
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func splitPath(p string) []string {
	var parts []string
	for _, seg := range split(p, '/') {
		if seg != "" {
			parts = append(parts, seg)
		}
	}
	return parts
}

func split(s string, sep rune) []string {
	var parts []string
	cur := ""
	for _, r := range s {
		if r == sep {
			parts = append(parts, cur)
			cur = ""
		} else {
			cur += string(r)
		}
	}
	parts = append(parts, cur)
	return parts
}

func PayeeGetOneAPI(store *PayeePostgresDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := splitPath(r.URL.Path)
		if len(parts) < 2 {
			http.Error(w, "id missing in path", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		p, err := store.GetByID(context.Background(), id)
		if err != nil {
			http.Error(w, "record not found", http.StatusNotFound)
			return
		}

		resp := PayeeGETResponse{
			ID:              p.id,
			BeneficiaryName: p.beneficiaryName,
			BeneficiaryCode: p.beneficiaryCode,
			AccNo:           p.accNo,
			IFSC:            p.ifsc,
			BankName:        p.bankName,
			Email:           p.email,
			Mobile:          p.mobile,
			PayeeCategory:   p.payeeCategory,
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
