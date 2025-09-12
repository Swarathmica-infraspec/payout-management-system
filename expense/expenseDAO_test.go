package expense

import (
	"context"
	"database/sql"
	payee "payoutmanagementsystem/payee"
	"testing"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal("database conncetion error")
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
		t.Fatal("db connection failed")
	}
	store := NewPostgresExpenseDB(db)
	defer clearExpenses(t, db)

	e, err := NewExpense("Lunch", 450.00, "2025-08-27", "Food", "Team lunch", 1, "/lunch.jpg")
	if err != nil {
		t.Fatalf("failed to create expense struct: %v", err)
	}

	id, err := store.Insert(context.Background(), e)
	if err != nil {
		t.Fatal("db connection failed")
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
func TestListExpensesForPayout(t *testing.T) {
	db := setupTestDB(t)
	expenseStore := NewPostgresExpenseDB(db)
	payeeStore := payee.PostgresPayeeDB(db)
	defer clearExpenses(t, db)

	payee1, err := payee.NewPayee("Abdef", "1901", 1934067090123856, "CBIN0123459", "CBI", "abdef@gmail.com", 9127960780, "Employee")
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	payeeId1, err := payeeStore.Insert(context.Background(), payee1)
	if err != nil {
		t.Fatalf("failed to insert payee: %v", err)
	}

	expense1, err := NewExpense("Lunch", 150.00, "2025-09-10", "food", "Team lunch", payeeId1, "/lunch.jpg")
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	expenseId1, err := expenseStore.Insert(context.Background(), expense1)
	if err != nil {
		t.Fatal("Insertion failed")
	}

	expense2, err := NewExpense("Taxi", 50.00, "2025-09-09", "travel", "Airport Drop", expenseId1, "/taxi.jpg")
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	expenseId2, err := expenseStore.Insert(context.Background(), expense2)
	if err != nil {
		t.Fatal("Insertion failed")
	}

	expense3, err := NewExpense("Taxi", 70.00, "2025-09-09", "travel", "Airport Drop", expenseId2, "/taxi.jpg")
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}
	expense3.status = "Paid"
	_, err = expenseStore.Insert(context.Background(), expense3)
	if err != nil {
		t.Fatal("Insertion failed")
	}

	expensesList, total, err := expenseStore.ListExpensesForPayout(context.Background())

	if err != nil {
		t.Fatal("Insertion failed for expense")
	}
	if len(expensesList) != 2 {
		t.Errorf("length of expenses has to be 2")
	}
	if total != 200.00 {
		t.Errorf("total has to be 200.00 but got %f", total)
	}
}
