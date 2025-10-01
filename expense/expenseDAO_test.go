package expense

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestInsertExpense(t *testing.T) {
	db := setupTestDB(t)
	defer func() { _ = db.Close() }()
	store := ExpenseDB(db)
	ctx := context.Background()

	e, err := NewExpense("Lunch", 450.00, futureDate(baseDate, 1), "Food", "Team lunch", 1, "/lunch.jpg")
	require.NoError(t, err, "failed to create expense struct")

	id, err := store.Insert(ctx, e)
	require.NoError(t, err, "insert operation failed")
	defer func() {
		_, _ = db.Exec("DELETE FROM expenses WHERE id = $1", id)
	}()

	var title, category, notes, receiptURI string
	var amount float64
	var payeeID int
	err = db.QueryRow(`
		SELECT title, amount, category, notes, payee_id, receipt_uri
		FROM expenses WHERE id = $1`, id).
		Scan(&title, &amount, &category, &notes, &payeeID, &receiptURI)
	require.NoError(t, err, "failed to query expense")

	assert.Equal(t, e.title, title)
	assert.Equal(t, e.amount, amount)
	assert.Equal(t, e.category, category)
	assert.Equal(t, e.notes, notes)
	assert.Equal(t, e.payeeID, payeeID)
	assert.Equal(t, e.receiptURI, receiptURI)
}

func TestGetExpenseByID(t *testing.T) {
	db := setupTestDB(t)
	defer func() { _ = db.Close() }()
	store := ExpenseDB(db)
	ctx := context.Background()

	var id int
	err := db.QueryRow(`
		INSERT INTO expenses (title, amount, date_incurred, category, notes, payee_id, receipt_uri)
		VALUES ('Dinner', 700.00, '2025-09-01', 'Food', 'Team dinner', 2, '/dinner.jpg')
		RETURNING id`).Scan(&id)
	require.NoError(t, err, "failed to insert fixture expense")
	defer func() {
		_, _ = db.Exec("DELETE FROM expenses WHERE id = $1", id)
	}()

	got, err := store.GetByID(ctx, id)
	require.NoError(t, err, "failed to fetch expense")

	title := "Dinner"
	amount := 700.00
	category := "Food"
	notes := "Team dinner"
	payeeID := 2
	receiptURI := "/dinner.jpg"

	assert.Equal(t, title, got.title)
	assert.Equal(t, amount, got.amount)
	assert.Equal(t, category, got.category)
	assert.Equal(t, notes, got.notes)
	assert.Equal(t, payeeID, got.payeeID)
	assert.Equal(t, receiptURI, got.receiptURI)
}
