package payoutmanagementsystem

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.POST("/expense", func(c *gin.Context) {
		var req struct {
			Title        string  `json:"title"`
			Amount       float64 `json:"amount"`
			DateIncurred string  `json:"dateIncurred"`
			Category     string  `json:"category"`
			Notes        string  `json:"notes"`
			PayeeID      int     `json:"payeeID"`
			ReceiptURI   string  `json:"receiptURI"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": 1})
	})

	return r
}

func TestExpensePostAPISuccess(t *testing.T) {
	router := setupRouter()

	payload := map[string]interface{}{
		"title":        "Lunch",
		"amount":       450.00,
		"dateIncurred": "2025-08-27",
		"category":     "Food",
		"notes":        "Team Lunch",
		"payeeID":      10,
		"receiptURI":   "/Desktop/lunch_bill.jpg",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/expense", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestExpensePostAPIInvalidJSON(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/expense", bytes.NewBufferString("{bad json}"))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
