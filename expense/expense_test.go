package expense

import (
	"testing"
)

var validExpenseTests = []struct {
	title        string
	amount       float64
	dateIncurred string
	category     string
	notes        string
	payeeID      int
	receiptURI   string
}{
	{"Lunch", 450.00, "2025-08-27", "Food", "Team lunch", 10, "/Desktop/lunch.jpg"},
	{"Travel", 120.00, "2025-08-26", "Transport", "Bus fare", 11, "/var/docs/paper-receipt.png"},
	{"Paper", 20, "2025-08-21", "Supplies", "For printer", 13, "/var/docs/paper-receipt.png"},
	{"Paper", 20, "2025-08-21", "Supplies", "For printer", 13, "/var/docs/paper-receipt.png"},
	{"Hotel", 2100, "2025-08-25", "Accommodation", "Stay", 1, "/var/docs/paper-receipt.png"},
}

func TestValidateExpenseWithValidValues(t *testing.T) {
	for _, tt := range validExpenseTests {
		_, err := NewExpense(tt.title, tt.amount, tt.dateIncurred, tt.category, tt.notes, tt.payeeID, tt.receiptURI)
		if err != nil {
			t.Fatalf("Expense can be created, but got: %v", err)
		}
	}
}

var invalidExpenseTests = []struct {
	title        string
	amount       float64
	dateIncurred string
	category     string
	notes        string
	payeeID      int
	receiptURI   string
	expectedErr  error
}{
	{"", 450.00, "2025-08-27", "Food", "Team lunch", 10, "https://receipts.com/lunch.jpg", ErrInvalidTitle},
	{"Travel", 0, "2025-08-27", "Travel", "Bus fare", 11, "", ErrInvalidAmount},
	{"Snacks", 55, "2025-08-32", "Food", "Evening snacks", 12, "", ErrInvalidDate},
	{"Snacks", 55, "2025-13-30", "Food", "Evening snacks", 12, "", ErrInvalidDate},
	{"Snacks", 55, "1999-12-24", "Food", "Evening snacks", 12, "", ErrInvalidDate},
	{"Paper", 20, "2025-08-21", "", "For printer", 13, "", ErrInvalidCategory},
	{"Hotel", 2100, "2025-08-25", "Accommodation", "Stay", -1, "", ErrInvalidPayeeID},
	{"Stationery", 200, "2025-08-24", "Office", "Pens", 14, "bill", ErrInvalidReceiptURI},
}

func TestValidateExpenseWithInvalidValues(t *testing.T) {
	for _, tt := range invalidExpenseTests {
		_, err := NewExpense(tt.title, tt.amount, tt.dateIncurred, tt.category, tt.notes, tt.payeeID, tt.receiptURI)
		if err != tt.expectedErr {
			t.Fatalf("Expected Error: %v but Actual Error: %v", tt.expectedErr, err)
		}
	}
}
