package payee

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type PayeeRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	AccNo    int    `json:"account_number"`
	IFSC     string `json:"ifsc"`
	Bank     string `json:"bank"`
	Email    string `json:"email"`
	Mobile   int    `json:"mobile"`
	Category string `json:"category"`
}

func PayeePostAPI(store PayeeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data PayeeRequest
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON body"}); err != nil {
				log.Printf("Failed to encode response (invalid JSON body): %v", err)
			}
			return
		}

		p, err := NewPayee(data.Name, data.Code, data.AccNo, data.IFSC, data.Bank, data.Email, data.Mobile, data.Category)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "Invalid payee data"}); err != nil {
				log.Printf("Failed to encode response (invalid payee data): %v", err)
			}
			return
		}

		id, err := store.Insert(context.Background(), p)
		if err != nil {

			var errMsg string
			var status int

			switch err {
			case ErrDuplicateCode:
				errMsg = "beneficiary code"
				status = http.StatusConflict
			case ErrDuplicateAccount:
				errMsg = "account number"
				status = http.StatusConflict
			case ErrDuplicateEmail:
				errMsg = "email"
				status = http.StatusConflict
			case ErrDuplicateMobile:
				errMsg = "mobile"
				status = http.StatusConflict
			default:
				errMsg = "internal server error"
				status = http.StatusInternalServerError
			}

			w.WriteHeader(status)
			if status == http.StatusConflict {
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Payee with the same " + errMsg + " already exists"}); err != nil {
					log.Printf("Failed to encode response (column value repetition): %v", err)
				}
			} else {
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong"}); err != nil {
					log.Printf("Failed to encode response (conflict due to server error): %v", err)
				}
				log.Printf("Internal error: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]any{"id": id}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
	}
}

func SetupRouter(store PayeeRepository) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/payees", PayeePostAPI(store))

	return mux
}
