package payout

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/payouts/preview", func(c *gin.Context) {
		expenses := []map[string]interface{}{
			{
				"expenseID":       1,
				"title":           "Lunch",
				"amount":          150.00,
				"dateIncurred":    "2025-09-10",
				"beneficiaryName": "Alice",
				"beneficiaryCode": "BEN001",
				"accountNumber":   12345678,
				"ifscCode":        "ABCD0123456",
				"bankName":        "HDFC Bank",
				"email":           "alice@example.com",
			},
		}

		total := 0.0
		for _, e := range expenses {
			total += e["amount"].(float64)
		}

		c.JSON(http.StatusOK, gin.H{
			"expenses": expenses,
			"total":    total,
		})
	})

	return r
}

func TestPayoutPreviewAPISuccess(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/payouts/preview", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

	body := w.Body.String()
	if !strings.Contains(body, "Lunch") {
		t.Errorf("expected response to contain Lunch, got %s", body)
	}
	if !strings.Contains(body, `"total":150`) {
		t.Errorf("expected response to contain total=150, got %s", body)
	}
}
