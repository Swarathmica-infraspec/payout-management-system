package payoutmanagementsystem

import (
	"context"
	"database/sql"
	"testing"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skip("skipping DB connection due to error:", err)
	}
	return db
}

func TestCreateAndGetExpense(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Skip("db connection not available, skipping test")
	}
	store := NewPostgresExpenseDB(db)

	e, err := NewExpense("Lunch", 450.00, "2025-08-27", "Food", "Team lunch", 1, "/lunch.jpg")
	if err != nil {
		t.Fatalf("failed to create expense struct: %v", err)
	}

	id, err := store.Insert(context.Background(), e)
	if err != nil {
		t.Skip("skipping insertion due to error:", err)
	}

	defer func() {
		if _, err := db.Exec("DELETE FROM expenses WHERE id = $1", id); err != nil {
			t.Errorf("failed to clean up expense id %d: %v", id, err)
		}
	}()

	got, err := store.GetByID(context.Background(), id)
	if err != nil {
		t.Fatalf("failed to fetch expense: %v", err)
	}

	if got.title != e.title {
		t.Errorf("expected title %q, got %q", e.title, got.title)
	}
	if got.amount != e.amount {
		t.Errorf("expected amount %v, got %v", e.amount, got.amount)
	}
	if got.dateIncurred != e.dateIncurred {
		t.Errorf("expected date %q, got %q", e.dateIncurred, got.dateIncurred)
	}
	if got.category != e.category {
		t.Errorf("expected category %q, got %q", e.category, got.category)
	}
	if got.notes != e.notes {
		t.Errorf("expected notes %q, got %q", e.notes, got.notes)
	}
	if got.payeeID != e.payeeID {
		t.Errorf("expected payeeID %d, got %d", e.payeeID, got.payeeID)
	}
	if got.receiptURI != e.receiptURI {
		t.Errorf("expected receiptURI %q, got %q", e.receiptURI, got.receiptURI)
	}
}
