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

func respondError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

func respondSuccess(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to encode success response: %v", err)
	}
}
func mapInsertError(err error) (int, string) {
	switch err {
	case ErrDuplicateCode:
		return http.StatusConflict, "Payee with the same beneficiary code already exists"
	case ErrDuplicateAccount:
		return http.StatusConflict, "Payee with the same account number already exists"
	case ErrDuplicateEmail:
		return http.StatusConflict, "Payee with the same email already exists"
	case ErrDuplicateMobile:
		return http.StatusConflict, "Payee with the same mobile already exists"
	default:
		return http.StatusInternalServerError, "Something went wrong"
	}
}

func handleInsertError(w http.ResponseWriter, err error) {
	status, message := mapInsertError(err)

	if status == http.StatusInternalServerError {
		log.Printf("Internal error: %v", err)
		respondError(w, status, "Something went wrong")
		return
	}

	respondError(w, status, message)
}

func PayeePostAPI(store PayeeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data PayeeRequest

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			respondError(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		p, err := NewPayee(data.Name, data.Code, data.AccNo, data.IFSC, data.Bank, data.Email, data.Mobile, data.Category)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid payee data")
		}

		id, err := store.Insert(context.Background(), p)
		if err != nil {
			handleInsertError(w, err)
			return
		}

		respondSuccess(w, http.StatusCreated, map[string]any{"id": id})
	}
}

func SetupRouter(store PayeeRepository) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/payees", PayeePostAPI(store))

	return mux
}
