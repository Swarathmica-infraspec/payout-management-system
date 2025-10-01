package payee

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
	require.NoError(t, err, "failed to connect to DB")
	err = db.Ping()
	require.NoError(t, err, "failed to ping DB")
	return db
}

func TestInsertPayee(t *testing.T) {
	db := setupTestDB(t)
	store := PayeeDB(db)
	ctx := context.Background()

	p, err := NewPayee("Abc", "136", 1234567890123456, "CBIN0123459", "CBI", "abc@gmail.com", 9123456780, "Employee")
	require.NoError(t, err, "failed to create payee")

	id, err := store.Insert(ctx, p)
	require.NoError(t, err, "failed to insert payee")

	defer func() {
		_, err := db.Exec("DELETE FROM payees WHERE id = $1", id)
		assert.NoError(t, err, "failed to clean up payee")
	}()

	var code, name, bank, ifsc, email, category string
	var accNo int
	var mobile int
	err = db.QueryRow(`
		SELECT beneficiary_code, beneficiary_name, account_number, ifsc_code, bank_name, email, mobile, payee_category
		FROM payees WHERE id = $1`, id).
		Scan(&code, &name, &accNo, &ifsc, &bank, &email, &mobile, &category)

	require.NoError(t, err, "failed to query payee")

	assert.Equal(t, p.beneficiaryCode, code)
	assert.Equal(t, p.beneficiaryName, name)
	assert.Equal(t, p.accNo, accNo)
	assert.Equal(t, p.ifsc, ifsc)
	assert.Equal(t, p.bankName, bank)
	assert.Equal(t, p.email, email)
	assert.Equal(t, p.mobile, mobile)
	assert.Equal(t, p.payeeCategory, category)
}
func TestInsertPayeeWithDuplicateValues(t *testing.T) {
	db := setupTestDB(t)
	store := PayeeDB(db)
	ctx := context.Background()

	// Insert the original payee
	original, err := NewPayee("Abc", "136", 1234567890123456, "CBIN0123459", "CBI", "abc@gmail.com", 9123456780, "Employee")
	require.NoError(t, err, "failed to create original payee")

	id, err := store.Insert(ctx, original)
	require.NoError(t, err, "failed to insert original payee")
	defer func() {
		_, err := db.Exec("DELETE FROM payees WHERE id = $1", id)
		assert.NoError(t, err, "failed to clean up original payee")
	}()

	tests := []struct {
		name      string
		duplicate func() *payee
		wantErr   error
	}{
		{
			name: "duplicate beneficiary code",
			duplicate: func() *payee {
				p, err := NewPayee("Abc", "136", 1234567800123456, "CBIN0123459", "CBI", "abcd@gmail.com", 9127456780, "Employee")
				require.NoError(t, err)
				return p
			},
			wantErr: ErrDuplicateCode,
		},
		{
			name: "duplicate account number",
			duplicate: func() *payee {
				p, err := NewPayee("Xyz", "137", 1234567890123456, "CBIN0123460", "CBI", "x@gmail.com", 9123456790, "Employee")
				require.NoError(t, err)
				return p
			},
			wantErr: ErrDuplicateAccount,
		},
		{
			name: "duplicate email",
			duplicate: func() *payee {
				p, err := NewPayee("Pqr", "138", 1234567890123450, "CBIN0123461", "CBI", "abc@gmail.com", 9123456800, "Employee")
				require.NoError(t, err)
				return p
			},
			wantErr: ErrDuplicateEmail,
		},
		{
			name: "duplicate mobile",
			duplicate: func() *payee {
				p, err := NewPayee("Xyz", "137", 9876543210987654, "CBIN0123460", "CBI", "xyz@gmail.com", 9123456780, "Employee")
				require.NoError(t, err)
				return p
			},
			wantErr: ErrDuplicateMobile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duplicatePayee := tt.duplicate()
			_, err := store.Insert(ctx, duplicatePayee)
			require.ErrorIs(t, err, tt.wantErr)
		})
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

	require.NoError(t, err, "failed to insert payee")
	defer func() {
		_, err := db.Exec("DELETE FROM payees WHERE id = $1", id)
		assert.NoError(t, err, "failed to clean up payee")
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

	require.NoError(t, err, "failed to fetch payee")

	assert.Equal(t, name, got.beneficiaryName)
	assert.Equal(t, code, got.beneficiaryCode)
	assert.Equal(t, accNo, got.accNo)
	assert.Equal(t, ifsc, got.ifsc)
	assert.Equal(t, bank, got.bankName)
	assert.Equal(t, email, got.email)
	assert.Equal(t, mobile, got.mobile)
	assert.Equal(t, category, got.payeeCategory)
}
