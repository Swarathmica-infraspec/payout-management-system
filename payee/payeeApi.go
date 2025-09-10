package payee

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

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

func PayeePostAPI(store *PayeePostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name     string `json:"name"`
			Code     string `json:"code"`
			AccNo    int    `json:"account_number"`
			IFSC     string `json:"ifsc"`
			Bank     string `json:"bank"`
			Email    string `json:"email"`
			Mobile   int    `json:"mobile"`
			Category string `json:"category"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
			return
		}

		p, err := NewPayee(req.Name, req.Code, req.AccNo, req.IFSC, req.Bank, req.Email, req.Mobile, req.Category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}

		id, err := store.Insert(context.Background(), p)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB insert failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

func PayeeGetApi(store *PayeePostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		payees, err := store.List(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed", "details": err.Error()})
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

		c.JSON(http.StatusOK, resp)
	}
}

func PayeeGetOneApi(store *PayeePostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		p, err := store.GetByID(context.Background(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
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

		c.JSON(http.StatusOK, resp)
	}
}

func PayeeUpdateApi(store *PayeePostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req struct {
			Name     string `json:"name"`
			Code     string `json:"code"`
			AccNo    int    `json:"account_number"`
			IFSC     string `json:"ifsc"`
			Bank     string `json:"bank"`
			Email    string `json:"email"`
			Mobile   int    `json:"mobile"`
			Category string `json:"category"`
		}

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
			return
		}

		p, err := NewPayee(req.Name, req.Code, req.AccNo, req.IFSC, req.Bank, req.Email, req.Mobile, req.Category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
			return
		}
		p.id = id

		_, err = store.Update(context.Background(), p)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB update failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	}
}

func PayeeDeleteApi(store *PayeePostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}

		err = store.Delete(context.Background(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB delete failed", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}

func SetupRouter(r *gin.RouterGroup, store *PayeePostgresDB) {
	r.POST("/payees", PayeePostAPI(store))
	r.GET("/payees", PayeeGetApi(store))
	r.GET("/payees/:id", PayeeGetOneApi(store))
	r.PUT("/payees/:id", PayeeUpdateApi(store))
	r.DELETE("/payees/:id", PayeeDeleteApi(store))

}
