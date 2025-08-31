package payoutmanagementsystem

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var store *PayeePostgresDB

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

func initStore() *PayeePostgresDB {
	if store != nil {
		return store
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	store = PostgresPayeeDB(db)
	return store
}

func PayeePostAPI(c *gin.Context) {

	store := initStore()

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

func PayeeGetApi(c *gin.Context) {
	store := initStore()

	payees, err := store.List(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed", "details": err.Error()})
		return
	}

	var resp []PayeeGETResponse
	for _, p := range payees {
		resp = append(resp, PayeeGETResponse{
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

func PayeeGetOneApi(c *gin.Context) {
	store := initStore()

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

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/payees", PayeePostAPI)
	r.GET("/payees", PayeeGetApi)
	r.GET("/payees/:id", PayeeGetOneApi)
	return r
}
