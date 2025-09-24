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

func TestInsertPayee(t *testing.T) {
	db := setupTestDB(t)
	store := PayeeDB(db)
	ctx := context.Background()

	p, _ := NewPayee("Abc", "136", 1234567890123456, "CBIN0123459", "CBI", "abc@gmail.com", 9123456780, "Employee")

	id, err := store.Insert(ctx, p)
	if err != nil {
		t.Fatalf("failed to insert payee: %v", err)
	}
	defer func() {
		if _, err := db.Exec("DELETE FROM payees WHERE id = $1", id); err != nil {
			t.Logf("failed to clear payee: %v", err)
		}
	}()

	var code, name, bank, ifsc, email, category string
	var accNo int
	var mobile int
	err = db.QueryRow(`
		SELECT beneficiary_code, beneficiary_name, account_number, ifsc_code, bank_name, email, mobile, payee_category
		FROM payees WHERE id = $1`, id).
		Scan(&code, &name, &accNo, &ifsc, &bank, &email, &mobile, &category)
	if err != nil {
		t.Fatalf("failed to query payee: %v", err)
	}

	if code != p.beneficiaryCode {
		t.Errorf("expected %+v, got %+v", code, p.beneficiaryCode)
	}

	if name != p.beneficiaryName {
		t.Errorf("expected %+v, got %+v", name, p.beneficiaryName)
	}
	if accNo != p.accNo {
		t.Errorf("expected %+v, got %+v", accNo, p.accNo)
	}

	if ifsc != p.ifsc {
		t.Errorf("expected %+v, got %+v", ifsc, p.ifsc)
	}
	if bank != p.bankName {
		t.Errorf("expected %+v, got %+v", bank, p.bankName)
	}
	if email != p.email {
		t.Errorf("expected %+v, got %+v", email, p.email)
	}
	if mobile != p.mobile {
		t.Errorf("expected %+v, got %+v", mobile, p.mobile)
	}
	if category != p.payeeCategory {
		t.Errorf("expected %+v, got %+v", mobile, p.mobile)
	}
}

func TestGetPayeeByID(t *testing.T) {
	db := setupTestDB(t)
	store := PayeeDB(db)
	ctx := context.Background()

	var id int
	err := db.QueryRow(`
		INSERT INTO payees (beneficiary_name, beneficiary_code, account_number, ifsc_code, bank_name, email, mobile, payee_category)
		VALUES ('Abc','136',1234567890123456,'CBIN0123459','CBI','abc@gmail.com',9123456780,'Employee')
		RETURNING id`).Scan(&id)
	if err != nil {
		t.Fatalf("failed to insert fixture payee: %v", err)
	}
	defer func() {
		if _, err := db.Exec("DELETE FROM payees WHERE id = $1", id); err != nil {
			t.Logf("failed to clear payee: %v", err)
		}
	}()

	got, err := store.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("failed to fetch payee: %v", err)
	}

	if got.beneficiaryName != "Abc" {
		t.Errorf("expected name Abc, got %s", got.beneficiaryName)
	}
}
