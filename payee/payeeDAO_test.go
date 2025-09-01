package payoutmanagementsystem

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to DB: %v", err)
		// t.Skip("skipping connection")
	}
	return db
}
func clearPayees(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE payees RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to clear table: %v", err)
	}
}

func TestInsertAndGetPayee(t *testing.T) {
	db := setupTestDB(t)
	store := PostgresPayeeDB(db)
	defer clearPayees(t, db)

	p, err := NewPayee("Abc", "123", 1234567890123456, "CBIN012345", "CBI", "abc@gmail.com", 9123456780, "Employee")
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	id, err := store.Insert(context.Background(), p)
	if err != nil {
		t.Fatalf("failed to insert payee: %v", err)
	}

	defer func() {
		if _, err := db.Exec("DELETE FROM payees WHERE id = $1", id); err != nil {
			t.Errorf("warning: failed to clean up payee id %d: %v", id, err)
		}
	}()

	got, err := store.GetByID(context.Background(), id)
	if err != nil {
		t.Fatalf("failed to fetch payee: %v", err)
	}

	if got.beneficiaryCode != p.beneficiaryCode {
		t.Errorf("expected beneficiary code: %s, got: %s", p.beneficiaryCode, got.beneficiaryCode)
	}
}

func TestListPayees(t *testing.T) {
	db := setupTestDB(t)
	store := PostgresPayeeDB(db)
	defer clearPayees(t, db)

	p, err := NewPayee("Xyz", "456", 1234567890123456, "HDFC000123", "HDFC", "xyz@gmail.com", 9876543210, "Vendor")
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	id, err := store.Insert(context.Background(), p)
	if err != nil {
		t.Skip("skipping insertion due to DB issue")
	}

	defer func() {
		if _, err := db.Exec("DELETE FROM payees WHERE id = $1", id); err != nil {
			t.Errorf("warning: failed to clean up payee id %d: %v", id, err)
		}
	}()

	_, err = store.List(context.Background())
	if err != nil {
		t.Fatalf("failed to list payees: %v", err)
		// t.Skip("skipping error check for List")
	}
}

func TestUpdatePayee(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer clearPayees(t, db)

	store := PostgresPayeeDB(db)

	p, _ := NewPayee("Abc", "123", 1234567890123456, "CBIN012345", "CBI", "abc@gmail.com", 9123456780, "Employee")
	id, err := store.Insert(ctx, p)
	if err != nil {
		t.Fatalf("failed to insert payee: %v", err)
	}

	originalPayee, _ := store.GetByID(ctx, id)

	updatedName := "cat"

	originalPayee.beneficiaryName = updatedName

	updated, err := store.Update(ctx, originalPayee)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if updated.beneficiaryName != updatedName {
		t.Errorf("expected name %q, got %q", updatedName, originalPayee.beneficiaryName)
	}
}
