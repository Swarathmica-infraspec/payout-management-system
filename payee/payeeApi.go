package payee

import (
	"context"
	"encoding/json"
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
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON body"})
			return
		}

		p, err := NewPayee(data.Name, data.Code, data.AccNo, data.IFSC, data.Bank, data.Email, data.Mobile, data.Category)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid payee data"})
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
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "Payee already exists with the same: " + errMsg})
			} else {
				_ = json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
	}
}

func SetupRouter(store PayeeRepository) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/payees", PayeePostAPI(store))

	return mux
}
