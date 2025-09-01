package payoutmanagementsystem

import (
	"context"
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var store *ExpensePostgresDB

type ExpenseGETResponse struct {
	Title        string `json:"title"`
	Amount       int    `json:"amount"`
	DateIncurred string `json:"dateIncurred"`
	Category     string `json:"category"`
	Notes        string `json:"notes"`
	PayeeID      int    `json:"payeeID"`
	ReceiptURI   string `json:"receiptURI"`
}

func initStore() *ExpensePostgresDB {
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
	store = NewPostgresExpenseDB(db)
	return store
}

func ExpensePostAPI(c *gin.Context) {

	store := initStore()

	var req struct {
		Title        string `json:"title"`
		Amount       int    `json:"amount"`
		DateIncurred string `json:"dateIncurred"`
		Category     string `json:"category"`
		Notes        string `json:"notes"`
		PayeeID      int    `json:"payeeID"`
		ReceiptURI   string `json:"receiptURI"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	p, err := NewExpense(req.Title, float64(req.Amount), req.DateIncurred, req.Category, req.Notes, req.PayeeID, req.ReceiptURI)
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

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/expense", ExpensePostAPI)
	return r
}
