package expense

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func futureDate(baseDate time.Time, daysAhead int) string {
	return baseDate.AddDate(0, 0, daysAhead).Format("2006-01-02")
}

func pastDate(baseDate time.Time, daysBack int) string {
	return baseDate.AddDate(0, 0, -daysBack).Format("2006-01-02")
}

var (
	baseDate = time.Now()
)

var invalidExpenseTests = []struct {
	testName     string
	title        string
	amount       float64
	dateIncurred string
	category     string
	notes        string
	payeeID      int
	receiptURI   string
	expectedErr  error
	errMsg       string
}{
	{"TestInvalidExpenseWithEmptyTitle", "", 450.00, futureDate(baseDate, 1), "Food", "Team lunch", 10, "https://receipts.com/lunch.jpg", ErrInvalidTitle, "expense title cannot be empty"},
	{"TestInvalidExpenseOfAmount0", "Travel", 0, futureDate(baseDate, 1), "Travel", "Bus fare", 11, "", ErrInvalidAmount, "expense amount must be greater than zero"},
	{"TestInvalidExpenseWithWrongDate", "Snacks", 55, "2025-08-32", "Food", "Evening snacks", 12, "", ErrInvalidDate, "invalid date format: date exceeds 31"},
	{"TestInvalidExpenseWithWrongMonth", "Snacks", 55, "2025-13-30", "Food", "Evening snacks", 12, "", ErrInvalidDate, "invalid date format: month exceeds 12"},
	{"TestInvalidExpenseWithPastDate", "Lunch", 100, pastDate(baseDate, 1), "Food", "Past expense", 10, "/path/receipt.jpg", ErrInvalidDate, "date incurred cannot be in the past"},
	{"TestInvalidExpenseWithWrongCategory", "Paper", 20, futureDate(baseDate, 1), "", "For printer", 13, "", ErrInvalidCategory, "expense category cannot be empty"},
	{"TestInvalidExpenseWithInvalidPayeeID", "Hotel", 2100, futureDate(baseDate, 1), "Accommodation", "Stay", -1, "", ErrInvalidPayeeID, "payee ID cannot be negative"},
	{"TestInvalidExpenseWithInvalidReceiptURI", "Stationery", 200, futureDate(baseDate, 1), "Office", "Pens", 14, "bill", ErrInvalidReceiptURI, "invalid receipt URI"},
}

func TestInvalidExpense(t *testing.T) {
	for _, tt := range invalidExpenseTests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := NewExpense(tt.title, tt.amount, tt.dateIncurred, tt.category, tt.notes, tt.payeeID, tt.receiptURI)
			assert.ErrorIs(t, err, tt.expectedErr, "Error Test Case: %v", tt.errMsg)

		})
	}
}

func TestValidExpense(t *testing.T) {
	title := "Lunch"
	amount := 450.00
	dateIncurred := futureDate(time.Now(), 1)
	category := "Food"
	notes := "Team lunch"
	payeeID := 10
	receiptURI := "/Desktop/lunch.jpg"

	e, err := NewExpense(title, amount, dateIncurred, category, notes, payeeID, receiptURI)
	require.NoError(t, err, "expense should be created but got error")

	assert.Equal(t, title, e.title, "title should match")
	assert.Equal(t, amount, e.amount, "amount should match")
	assert.Equal(t, dateIncurred, e.dateIncurred, "dateIncurred should match")
	assert.Equal(t, category, e.category, "category should match")
	assert.Equal(t, notes, e.notes, "notes should match")
	assert.Equal(t, payeeID, e.payeeID, "payeeID should match")
	assert.Equal(t, receiptURI, e.receiptURI, "receiptURI should match")

}
