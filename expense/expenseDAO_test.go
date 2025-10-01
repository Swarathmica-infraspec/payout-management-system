package expense

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("TEST_DATABASE_URL")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
	}
	return db
}

func TestInsertAndGetPayee(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close DB connection: %v", err)
		}
	}()
	store := ExpenseDB(db)
	ctx := context.Background()

    e, err := NewExpense("Lunch", 450.00, "2025-08-27", "Food", "Team lunch", 1, "/lunch.jpg")
    if err != nil {
        t.Fatalf("failed to create expense struct: %v", err)
    }

    id, err := store.Insert(ctx, e)
    if err != nil {
        t.Fatalf("insert operation failed: %v", err)
    }

    defer func() {
        if _, err := db.Exec("DELETE FROM expenses WHERE id = $1", id); err != nil {
            t.Errorf("failed to clean up expense id %d: %v", id, err)
        }
    }()

    got, err := store.GetByID(ctx, id)
    if err != nil {
        t.Fatalf("failed to fetch expense: %v", err)
    }

    if got.title != e.title {
        t.Errorf("expected title %q, got %q", e.title, got.title)
    }
    if got.amount != e.amount {
        t.Errorf("expected amount %v, got %v", e.amount, got.amount)
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