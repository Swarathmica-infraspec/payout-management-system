package payee

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

func clearPayees(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE payees RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to clear table: %v", err)
	}
}


func TestInsertAndGetPayee(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close DB connection: %v", err)
		}
	}()
	store := PayeeDB(db)
	ctx := context.Background()

	p, err := NewPayee("Abc", "136", 1234567890123456, "CBIN0123459", "CBI", "abc@gmail.com", 9123456780, "Employee")
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	id, err := store.Insert(ctx, p)
	if err != nil {
		t.Fatalf("failed to insert payee: %v", err)
	}

	defer func() {
		if _, err := db.Exec("DELETE FROM payees WHERE id = $1", id); err != nil {
			t.Errorf("warning: failed to clean up payee id %d: %v", id, err)
		}
	}()

	got, err := store.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("failed to fetch payee: %v", err)
	}

	if got.beneficiaryCode != p.beneficiaryCode {
		t.Errorf("expected beneficiary code: %s, got: %s", p.beneficiaryCode, got.beneficiaryCode)
	}
	if got.beneficiaryName != p.beneficiaryName {
		t.Errorf("expected beneficiary name: %s, got: %s", p.beneficiaryName, got.beneficiaryName)
	}
	if got.accNo != p.accNo {
		t.Errorf("expected accNo: %d, got: %d", p.accNo, got.accNo)
	}
	if got.ifsc != p.ifsc {
		t.Errorf("expected IFSC code: %s, got: %s", p.ifsc, got.ifsc)
	}
	if got.bankName != p.bankName {
		t.Errorf("expected bank name: %s, got: %s", p.bankName, got.bankName)
	}
	if got.email != p.email {
		t.Errorf("expected email: %s, got: %s", p.email, got.email)
	}
	if got.mobile != p.mobile {
		t.Errorf("expected mobile: %d, got: %d", p.mobile, got.mobile)
	}
	if got.payeeCategory != p.payeeCategory {
		t.Errorf("expected beneficiary code: %s, got: %s", p.payeeCategory, got.payeeCategory)
	}
}

func TestListPayees(t *testing.T) {
	db := setupTestDB(t)
	store := PayeeDB(db)
	defer clearPayees(t, db)

	p, err := NewPayee("Xyz", "456", 1234567890123456, "HDFC0001213", "HDFC", "xyz@gmail.com", 9876543210, "Vendor")
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	id, err := store.Insert(context.Background(), p)
	if err != nil {
		t.Fatal("Insertion failed")
	}

	defer func() {
		if _, err := db.Exec("DELETE FROM payees WHERE id = $1", id); err != nil {
			t.Errorf("warning: failed to clean up payee id %d: %v", id, err)
		}
	}()

	_, err = store.List(context.Background())
	if err != nil {
		t.Fatalf("failed to list payees: %v", err)
	}
}
