package expense

import (
	"testing"
	"time"
)

func futureDate(daysAhead int) string {
	return time.Now().AddDate(0, 0, daysAhead).Format("2006-01-02")
}

func pastDate(daysBack int) string {
	return time.Now().AddDate(0, 0, -daysBack).Format("2006-01-02")
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

func TestInvalidExpense(t *testing.T) {
	for _, tt := range invalidExpenseTests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := NewExpense(tt.title, tt.amount, tt.dateIncurred, tt.category, tt.notes, tt.payeeID, tt.receiptURI)
			if err != tt.expectedErr {
				t.Fatalf("Expected Error: %v but Actual Error: %v", tt.expectedErr, err)
			}
		})
	}
}

func TestValidExpense(t *testing.T) {
	title := "Lunch"
	amount := 450.00
	dateIncurred := futureDate(1)
	category := "Food"
	notes := "Team lunch"
	payeeID := 10
	receiptURI := "/Desktop/lunch.jpg"

	e, err := NewExpense(title, amount, dateIncurred, category, notes, payeeID, receiptURI)
	if err != nil {
		t.Fatalf("expense should be created but got error: %v", err)
	}

	if e.title != title {
		t.Errorf("expected title: %v but got: %v", title, e.title)
	}

	if e.amount != amount {
		t.Errorf("expected amount: %v but got: %v", amount, e.amount)
	}

	if e.dateIncurred != dateIncurred {
		t.Errorf("expected date: %v but got: %v", dateIncurred, e.dateIncurred)
	}

	if e.category != category {
		t.Errorf("expected category: %v but got: %v", category, e.category)
	}

	if e.notes != notes {
		t.Errorf("expected notes: %v but got: %v", notes, e.notes)
	}

	if e.payeeID != payeeID {
		t.Errorf("expected payeeID: %v but got: %v", payeeID, e.payeeID)
	}

	if e.receiptURI != receiptURI {
		t.Errorf("expected receiptURI: %v but got: %v", receiptURI, e.receiptURI)
	}
}