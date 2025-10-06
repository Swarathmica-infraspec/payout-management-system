package payee

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("TEST_DATABASE_URL")

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err, "failed to connect to DB")
	err = db.Ping()
	require.NoError(t, err, "failed to ping DB")
	return db
}

func clearPayees(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE payees RESTART IDENTITY CASCADE")
	require.NoError(t, err, "failed to clear DB")
}

func TestInsertPayee(t *testing.T) {
	db := setupTestDB(t)
	store := PayeeDB(db)
	ctx := context.Background()

	defer clearPayees(t, db)

	p, err := NewPayee("Abc", "136", 1234567890123456, "CBIN0123459", "CBI", "abc@gmail.com", 9123456780, "Employee")
	require.NoError(t, err, "failed to create payee")

	id, err := store.Insert(ctx, p)
	require.NoError(t, err, "failed to insert payee")

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

	defer clearPayees(t, db)

	original, err := NewPayee("Abc", "136", 1234567890123456, "CBIN0123459", "CBI", "abc@gmail.com", 9123456780, "Employee")
	require.NoError(t, err, "failed to create original payee")

	_, err = store.Insert(ctx, original)
	require.NoError(t, err, "failed to insert original payee")

	tests := []struct {
		testName string
		nameArg  string
		code     string
		accNo    int
		ifsc     string
		bank     string
		email    string
		mobile   int
		category string
		wantErr  error
	}{
		{
			testName: "duplicate beneficiary code",
			nameArg:  "Abc",
			code:     "136",
			accNo:    1234567800123456,
			ifsc:     "CBIN0123459",
			bank:     "CBI",
			email:    "abcd@gmail.com",
			mobile:   9127456780,
			category: "Employee",
			wantErr:  ErrDuplicateCode,
		},
		{
			testName: "duplicate account number",
			nameArg:  "Xyz",
			code:     "137",
			accNo:    1234567890123456,
			ifsc:     "CBIN0123460",
			bank:     "CBI",
			email:    "x@gmail.com",
			mobile:   9123456790,
			category: "Employee",
			wantErr:  ErrDuplicateAccount,
		},
		{
			testName: "duplicate email",
			nameArg:  "Pqr",
			code:     "138",
			accNo:    1234567890123450,
			ifsc:     "CBIN0123461",
			bank:     "CBI",
			email:    "abc@gmail.com",
			mobile:   9123456800,
			category: "Employee",
			wantErr:  ErrDuplicateEmail,
		},
		{
			testName: "duplicate mobile",
			nameArg:  "Xyz",
			code:     "137",
			accNo:    9876543210987654,
			ifsc:     "CBIN0123460",
			bank:     "CBI",
			email:    "xyz@gmail.com",
			mobile:   9123456780,
			category: "Employee",
			wantErr:  ErrDuplicateMobile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			dup, err := NewPayee(tt.nameArg, tt.code, tt.accNo, tt.ifsc, tt.bank, tt.email, tt.mobile, tt.category)
			require.NoError(t, err)

			_, err = store.Insert(ctx, dup)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGetPayeeByID(t *testing.T) {
	db := setupTestDB(t)
	store := PayeeDB(db)
	ctx := context.Background()

	defer clearPayees(t, db)

	var id int
	err := db.QueryRow(`
		INSERT INTO payees (beneficiary_name, beneficiary_code, account_number, ifsc_code, bank_name, email, mobile, payee_category)
		VALUES ('Abc','136',1234567890123456,'CBIN0123459','CBI','abc@gmail.com',9123456780,'Employee')
		RETURNING id`).Scan(&id)

	require.NoError(t, err, "failed to insert payee")

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
func TestListPayees(t *testing.T) {
	db := setupTestDB(t)
	store := PayeeDB(db)
	defer clearPayees(t, db)

	p, err := NewPayee("Xyz", "456", 1234567890123456, "HDFC0001213", "HDFC", "xyz@gmail.com", 9876543210, "Vendor")
	require.NoError(t, err, "validation failed")

	_, err = store.Insert(context.Background(), p)
	require.NoError(t, err, "Insertion failed")

	payees, err := store.List(context.Background())
	require.NoError(t, err, "failed to list payees")

	assert.NotEmpty(t, payees, "expected at least one payee")
}

func TestUpdatePayee(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer clearPayees(t, db)

	store := PayeeDB(db)

	p, _ := NewPayee("Abc", "123", 1234567890123456, "CBIN0124345", "CBI", "abc@gmail.com", 9123456780, "Employee")
	id, err := store.Insert(ctx, p)
	require.NoError(t, err, "Insertion failed")

	originalPayee, _ := store.GetByID(ctx, id)

	updatedName := "cat"

	originalPayee.beneficiaryName = updatedName

	updated, err := store.Update(ctx, originalPayee)
	require.NoError(t, err, "Update failed")

	if updated.beneficiaryName != updatedName {
		t.Errorf("expected name %q, got %q", updatedName, updated.beneficiaryName)
	}
	assert.Equal(t, updatedName, updated.beneficiaryName)
}

func TestUpdatePayeeWithDuplicateValues(t *testing.T) {
	db := setupTestDB(t)
	store := PayeeDB(db)
	ctx := context.Background()
	defer clearPayees(t, db)

	p1, err := NewPayee("Abc", "111", 1234567890123456, "CBIN0001", "CBI", "abc@gmail.com", 9000000001, "Employee")
	require.NoError(t, err)
	_, err = store.Insert(ctx, p1)
	require.NoError(t, err)

	p2, err := NewPayee("Bravo", "222", 6543210987654321, "HDFC0002", "HDFC", "bravo@gmail.com", 9000000002, "Vendor")
	require.NoError(t, err)
	id2, err := store.Insert(ctx, p2)
	require.NoError(t, err)

	tests := []struct {
		testName string
		targetID int
		updateFn func(p *payee)
		wantErr  error
	}{
		{
			testName: "duplicate beneficiary code",
			targetID: id2,
			updateFn: func(p *payee) {
				p.beneficiaryCode = "111"
			},
			wantErr: ErrDuplicateCode,
		},
		{
			testName: "duplicate account number",
			targetID: id2,
			updateFn: func(p *payee) {
				p.accNo = 1234567890123456
			},
			wantErr: ErrDuplicateAccount,
		},
		{
			testName: "duplicate email",
			targetID: id2,
			updateFn: func(p *payee) {
				p.email = "alpha@gmail.com" 
			},
			wantErr: ErrDuplicateEmail,
		},
		{
			testName: "duplicate mobile",
			targetID: id2,
			updateFn: func(p *payee) {
				p.mobile = 9000000001 
			},
			wantErr: ErrDuplicateMobile,
		},
		
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			original, err := store.GetByID(ctx, tt.targetID)
			require.NoError(t, err)

			tt.updateFn(original)

			_, err = store.Update(ctx, original)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
