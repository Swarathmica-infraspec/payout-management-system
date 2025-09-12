package payout

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"payoutmanagementsystem/expense"
)

type PayoutPreviewResponse struct {
	Expenses []expense.ExpenseWithPayee `json:"expenses"`
	Total    float64                    `json:"total"`
}

func PayoutPreviewAPI(store *expense.ExpensePostgresDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		expenses, total, err := store.ListExpensesForPayout(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to fetch payout preview",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, PayoutPreviewResponse{
			Expenses: expenses,
			Total:    total,
		})
	}
}

func SetupPayoutRoutes(r *gin.RouterGroup, store *expense.ExpensePostgresDB) {
	r.GET("/payouts/preview", PayoutPreviewAPI(store))
}
