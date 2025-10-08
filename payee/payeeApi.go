package payee

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		log.Printf("Failed to encode error response: %v", err)
	}
}

func respondSuccess(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
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

func parseFilterList(r *http.Request) FilterList {
	query := r.URL.Query()
	limit := 10
	offset := 0

	if l := query.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			if parsed < 1 {
				limit = 10
			} else if parsed > 100 {
				limit = 100
			} else {
				limit = parsed
			}
		}
	}

	if o := query.Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			if parsed < 0 {
				offset = 0
			} else {
				offset = parsed
			}
		}
	}

	sortBy := query.Get("sort_by")
	allowedSortFields := map[string]bool{
		"name": true, "bank": true, "category": true, "id": true,
	}
	if sortBy != "" && !allowedSortFields[sortBy] {
		sortBy = ""
	}

	sortOrder := query.Get("sort_order")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	return FilterList{
		Name:      query.Get("name"),
		Category:  query.Get("category"),
		Bank:      query.Get("bank"),
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Limit:     limit,
		Offset:    offset,
	}
}

func payeesToGETResponses(payees []payee) []PayeeGETResponse {
	resp := make([]PayeeGETResponse, 0, len(payees))
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
	return resp
}

func PayeeGetAPI(store PayeeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts := parseFilterList(r)

		payees, err := store.List(context.Background(), opts)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "DB query failed")
			return
		}

		respondSuccess(w, http.StatusOK, payeesToGETResponses(payees))
	}
}

func SetupRouter(store PayeeRepository) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/payees", PayeePostAPI(store))
	mux.HandleFunc("/payees/list", PayeeGetAPI(store))

	return mux
}
