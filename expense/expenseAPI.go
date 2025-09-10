package expense

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type ExpenseGETResponse struct {
	Title        string `json:"title"`
	Amount       int    `json:"amount"`
	DateIncurred string `json:"dateIncurred"`
	Category     string `json:"category"`
	Notes        string `json:"notes"`
	PayeeID      int    `json:"payeeID"`
	ReceiptURI   string `json:"receiptURI"`
}

func ExpensePostAPI(store *ExpensePostgresDB) gin.HandlerFunc {
	
	return func(c *gin.Context) {
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
}

func ExpenseGetApi(store *ExpensePostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		expenses, err := store.List(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed", "details": err.Error()})
			return
		}

		var resp []ExpenseGETResponse
		for _, e := range expenses {
			resp = append(resp, ExpenseGETResponse{
				Title: e.title,
				Amount: int(e.amount),
				DateIncurred: e.dateIncurred,
				Category: e.category,
				Notes: e.notes,
				PayeeID: e.payeeID,
				ReceiptURI: e.receiptURI,
			})
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ExpenseGetOneApi(store *ExpensePostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		e, err := store.GetByID(context.Background(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "record not found"})
			return
		}

		resp := ExpenseGETResponse{
			Title: e.title,
			Amount: int(e.amount),
			DateIncurred: e.dateIncurred,
			Category: e.category,
			Notes: e.notes,
			PayeeID: e.payeeID,
			ReceiptURI: e.receiptURI,
		}

		c.JSON(http.StatusOK, resp)
	}
}

func SetupRouter(r *gin.RouterGroup,store *ExpensePostgresDB) {
	r.POST("/expense", ExpensePostAPI(store))
	r.GET("/expense", ExpenseGetApi(store))
	r.GET("/expense/:id", ExpenseGetOneApi(store))
}
