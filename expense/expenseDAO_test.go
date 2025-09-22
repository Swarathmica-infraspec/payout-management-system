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
		t.Fatal("database connection error")
	}
	return db
}

func clearExpenses(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE expenses RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to clear table: %v", err)
	}
}

func TestInsertAndGetExpense(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		t.Fatal("database connection error")
	}
	store := NewPostgresExpenseDB(db)
	defer clearExpenses(t, db)

	e, err := NewExpense("Lunch", 450.00, "2025-08-27", "Food", "Team lunch", 1, "/lunch.jpg")
	if err != nil {
		t.Fatalf("failed to create expense struct: %v", err)
	}

	id, err := store.Insert(context.Background(), e)
	if err != nil {
		t.Fatalf("insert operation failed: %v", err)
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

func TestListExpenses(t *testing.T) {
	db := setupTestDB(t)
	store := NewPostgresExpenseDB(db)
	defer clearExpenses(t, db)

	p, err := NewExpense("Lunch", 450.00, "2025-08-27", "Food", "Team lunch", 1, "/lunch.jpg")
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	id, err := store.Insert(context.Background(), p)
	if err != nil {
		t.Fatal("Insertion failed")
	}

	defer func() {
		if _, err := db.Exec("DELETE FROM expenses WHERE payee_id = $1", id); err != nil {
			t.Errorf("warning: failed to clean up payee id %d: %v", id, err)
		}
	}()

	_, err = store.List(context.Background())
	if err != nil {
		t.Fatalf("failed to list payees: %v", err)
	}
}
