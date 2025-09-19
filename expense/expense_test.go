package expense

import (
	"fmt"
	"testing"
	"time"
)

func futureDate(daysAhead int) string {
	return time.Now().AddDate(0, 0, daysAhead).Format("2006-01-02")
}

func pastDate(daysBack int) string {
	return time.Now().AddDate(0, 0, -daysBack).Format("2006-01-02")
}

var validExpenseTests = []struct {
	title        string
	amount       float64
	dateIncurred string
	category     string
	notes        string
	payeeID      int
	receiptURI   string
}{
	{"Lunch", 450.00, futureDate(1), "Food", "Team lunch", 10, "/Desktop/lunch.jpg"},
	{"Travel", 120.00, futureDate(2), "Transport", "Bus fare", 11, "/var/docs/paper-receipt.png"},
	{"Paper", 20, futureDate(3), "Supplies", "For printer", 13, "/var/docs/paper-receipt.png"},
	{"Paper", 20, futureDate(4), "Supplies", "For printer", 13, "/var/docs/paper-receipt.png"},
	{"Hotel", 2100, futureDate(5), "Accommodation", "Stay", 1, "/var/docs/paper-receipt.png"},
}

func TestValidateExpenseWithValidValues(t *testing.T) {
	for i, tt := range validExpenseTests {
		t.Run(fmt.Sprintf("ValidExpenseCase%d_%s", i, tt.title), func(t *testing.T) {
			_, err := NewExpense(tt.title, tt.amount, tt.dateIncurred, tt.category, tt.notes, tt.payeeID, tt.receiptURI)
			if err != nil {
				t.Fatalf("Expense can be created, but got: %v", err)
			}
		})
	}
}

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
}{
	{"TestInvalidExpenseWithEmptyTitle", "", 450.00, futureDate(1), "Food", "Team lunch", 10, "https://receipts.com/lunch.jpg", ErrInvalidTitle},
	{"TestInvalidExpenseOfAmount0", "Travel", 0, futureDate(1), "Travel", "Bus fare", 11, "", ErrInvalidAmount},
	{"TestInvalidExpenseWithWrongDate", "Snacks", 55, "2025-08-32", "Food", "Evening snacks", 12, "", ErrInvalidDate},
	{"TestInvalidExpenseWithWrongMonth", "Snacks", 55, "2025-13-30", "Food", "Evening snacks", 12, "", ErrInvalidDate},
	{"TestInvalidExpenseWithYearBefore2025", "Snacks", 55, "1999-12-24", "Food", "Evening snacks", 12, "", ErrInvalidDate},
	{"TestInvalidExpenseWithPastDate", "Lunch", 100, pastDate(1), "Food", "Past expense", 10, "/path/receipt.jpg", ErrInvalidDate},
	{"TestInvalidExpenseWithWrongCategory", "Paper", 20, futureDate(1), "", "For printer", 13, "", ErrInvalidCategory},
	{"TestInvalidExpenseWithInvalidPayeeID", "Hotel", 2100, futureDate(1), "Accommodation", "Stay", -1, "", ErrInvalidPayeeID},
	{"TestInvalidExpenseWithInvalidReceiptURI", "Stationery", 200, futureDate(1), "Office", "Pens", 14, "bill", ErrInvalidReceiptURI},
}

func TestValidateExpenseWithInvalidValues(t *testing.T) {
	for _, tt := range invalidExpenseTests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := NewExpense(tt.title, tt.amount, tt.dateIncurred, tt.category, tt.notes, tt.payeeID, tt.receiptURI)
			if err != tt.expectedErr {
				t.Fatalf("Expected Error: %v but Actual Error: %v", tt.expectedErr, err)
			}
		})
	}
}
