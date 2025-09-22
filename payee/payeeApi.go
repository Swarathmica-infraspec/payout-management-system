package payee

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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

func PayeePostAPI(store PayeeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var data PayeeRequest
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		p, err := NewPayee(data.Name, data.Code, data.AccNo, data.IFSC, data.Bank, data.Email, data.Mobile, data.Category)
		if err != nil {
			http.Error(w, "Invalid payee data", http.StatusBadRequest)
			return
		}

		id, err := store.Insert(context.Background(), p)
		if err != nil {
			http.Error(w, "DB insertion failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": id})
	}
}

func PayeeGetAPI(store PayeeRepository) http.HandlerFunc {
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func PayeeGetOneAPI(store PayeeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := strings.TrimPrefix(r.URL.Path, "/payees/")
		id, err := strconv.Atoi(idStr)
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func PayeeUpdateAPI(store PayeeRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/payees/update/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid ID", http.StatusBadRequest)
			return
		}

		var req PayeeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		p, err := NewPayee(req.Name, req.Code, req.AccNo, req.IFSC, req.Bank, req.Email, req.Mobile, req.Category)
		if err != nil {
			http.Error(w, "validation failed: "+err.Error(), http.StatusBadRequest)
			return
		}
		p.id = id

		_, err = store.Update(context.Background(), p)
		if err != nil {
			http.Error(w, "DB update failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	}
}
