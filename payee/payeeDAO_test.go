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
	name := "Abc"
	code := "136"
	accNo := 1234567890123456
	ifsc := "CBIN0123459"
	bank := "CBI"
	email := "abc@gmail.com"
	mobile := 9123456780
	category := "Employee"
	if err != nil {
		t.Fatalf("failed to fetch payee: %v", err)
	}

	if got.beneficiaryName != name {
		t.Errorf("expected name %s, got %s", name, got.beneficiaryName)
	}
	if got.beneficiaryCode != code {
		t.Errorf("expected  code %s, got %s", code, got.beneficiaryCode)
	}
	if got.accNo != accNo {
		t.Errorf("expected accNo %d, got %d", accNo, got.accNo)
	}
	if got.ifsc != ifsc {
		t.Errorf("expected ifsc %s, got %s", ifsc, got.ifsc)
	}
	if got.bankName != bank {
		t.Errorf("expected bank %s, got %s", bank, got.bankName)
	}
	if got.email != email {
		t.Errorf("expected email %s, got %s", email, got.email)
	}
	if got.mobile != mobile {
		t.Errorf("expected mobile %d, got %d", mobile, got.mobile)
	}
	if got.payeeCategory != category {
		t.Errorf("expected category %s, got %s", category, got.payeeCategory)
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
