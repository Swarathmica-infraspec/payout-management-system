package payee

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.GET("/payees", func(c *gin.Context) {
		payees := []map[string]interface{}{
			{"id": 1, "name": "Alice"},
		}
		c.JSON(http.StatusOK, payees)
	})

	r.POST("/payees", func(c *gin.Context) {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": 1})
	})

	return r
}

func TestPayeePostAPISuccess(t *testing.T) {
	router := setupRouter()

	payload := map[string]interface{}{
		"name":           "Abc",
		"code":           "123",
		"account_number": 123456789,
		"ifsc":           "CBIN012345",
		"bank":           "CBI",
		"email":          "abc@example.com",
		"mobile":         9876543210,
		"category":       "Employee",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/payees", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestPayeePostAPIInvalidJSON(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("POST", "/payees", bytes.NewBufferString("{bad json}"))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPayeeGetAPISuccess(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/payees", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}

}
