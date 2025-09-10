package expense

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func ExpensePostAPI() {

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

	store := NewPostgresExpenseDB(db)

	router := gin.Default()

	router.POST("/expense", func(c *gin.Context) {
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
	})
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
