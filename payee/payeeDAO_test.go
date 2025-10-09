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
			code:     "139",
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

func TestListPayees(t *testing.T) {
	db := setupTestDB(t)
	store := PayeeDB(db)
	defer clearPayees(t, db)

	p1, _ := NewPayee("Alice", "A001", 1114567891234567, "HDFC0012345", "HDFC", "a@example.com", 9000000001, "Vendor")
	p2, _ := NewPayee("Bob", "B001", 2223456789012345, "SBIN0023478", "SBI", "b@example.com", 9000000002, "Employee")
	p3, _ := NewPayee("Charlie", "C001", 3334567890123456, "HDFC0033456", "HDFC", "c@example.com", 9000000003, "Vendor")

	_, _ = store.Insert(context.Background(), p1)
	_, _ = store.Insert(context.Background(), p2)
	_, _ = store.Insert(context.Background(), p3)

	tests := []struct {
		name       string
		opts       FilterList
		wantNames  []string
		wantIDs    []int
		minResults int
	}{
		{
			name:       "list all payees",
			opts:       FilterList{},
			minResults: 1,
		},
		{
			name:      "filter by name",
			opts:      FilterList{Name: "Alice"},
			wantNames: []string{"Alice"},
		},
		{
			name:      "filter by category",
			opts:      FilterList{Category: "Vendor"},
			wantNames: []string{"Alice", "Charlie"},
		},
		{
			name:      "filter by bank",
			opts:      FilterList{Bank: "SBI"},
			wantNames: []string{"Bob"},
		},
		{
			name:      "sort by id ASC",
			opts:      FilterList{SortBy: "id", SortOrder: "ASC"},
			wantNames: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name:      "sort by name ASC",
			opts:      FilterList{SortBy: "name", SortOrder: "ASC"},
			wantNames: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name:      "sort by name DESC",
			opts:      FilterList{SortBy: "name", SortOrder: "DESC"},
			wantNames: []string{"Charlie", "Bob", "Alice"},
		},
		{
			name:    "pagination: limit 1 offset 0 (first payee)",
			opts:    FilterList{SortBy: "id", SortOrder: "ASC", Limit: 1, Offset: 0},
			wantIDs: []int{1},
		},
		{
			name:    "pagination: limit 1 offset 1 (second payee)",
			opts:    FilterList{SortBy: "id", SortOrder: "ASC", Limit: 1, Offset: 1},
			wantIDs: []int{2},
		},
		{
			name:    "pagination: limit 2 offset 1 (second and third payees)",
			opts:    FilterList{SortBy: "id", SortOrder: "ASC", Limit: 2, Offset: 1},
			wantIDs: []int{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := store.List(context.Background(), tt.opts)
			require.NoError(t, err)

			if tt.minResults > 0 {
				assert.GreaterOrEqual(t, len(got), tt.minResults)
			}

			if len(tt.wantNames) > 0 {
				var gotNames []string
				for _, p := range got {
					gotNames = append(gotNames, p.beneficiaryName)
				}
				assert.Equal(t, tt.wantNames, gotNames)
			}

			if len(tt.wantIDs) > 0 {
				var gotIDs []int
				for _, p := range got {
					gotIDs = append(gotIDs, p.id)
				}
				assert.Equal(t, tt.wantIDs, gotIDs)
			}
		})
	}
}
func TestUpdatePayee(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer clearPayees(t, db)

	store := PayeeDB(db)

	p, err := NewPayee("Abc", "123", 1234567890123456, "CBIN0124345", "CBI", "abc@gmail.com", 9123456780, "Employee")
	require.NoError(t, err, "failed to create payee for update test")

	id, err := store.Insert(ctx, p)
	require.NoError(t, err, "Insertion failed")

	originalPayee := &payee{}

	err = db.QueryRowContext(ctx, `
		SELECT id, beneficiary_code, beneficiary_name, account_number, ifsc_code, bank_name, email, mobile, payee_category
		FROM payees WHERE id = $1`, id).
		Scan(&originalPayee.id, &originalPayee.beneficiaryCode, &originalPayee.beneficiaryName,
			&originalPayee.accNo, &originalPayee.ifsc, &originalPayee.bankName,
			&originalPayee.email, &originalPayee.mobile, &originalPayee.payeeCategory)

	require.NoError(t, err, "failed to fetch original payee")

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

	p1, err := NewPayee("Abc", "111", 1234567890123456, "CBIN0001234", "CBI", "abc@gmail.com", 9000000001, "Employee")
	require.NoError(t, err)
	_, err = store.Insert(ctx, p1)
	require.NoError(t, err)

	p2, err := NewPayee("Bcd", "222", 6543210987654321, "HDFC0001223", "HDFC", "bcd@gmail.com", 9000000002, "Vendor")
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
				p.email = "abc@gmail.com"
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

			original := &payee{}

			err := db.QueryRowContext(ctx, `
				SELECT id, beneficiary_code, beneficiary_name, account_number, ifsc_code, bank_name, email, mobile, payee_category
				FROM payees WHERE id = $1`, tt.targetID).
				Scan(&original.id, &original.beneficiaryCode, &original.beneficiaryName,
					&original.accNo, &original.ifsc, &original.bankName,
					&original.email, &original.mobile, &original.payeeCategory)

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
func TestSoftDeletePayee(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer clearPayees(t, db)

	store := PayeeDB(db)

	p, _ := NewPayee("Abc", "123", 1234567890123456, "CBIN0123456", "CBI", "abc@gmail.com", 9123456780, "Employee")
	id, err := store.Insert(ctx, p)
	require.NoError(t, err, "failed to insert payee")

	err = store.SoftDelete(ctx, id)
	require.NoError(t, err, "soft delete failed")

	got := &payee{}
	err = db.QueryRowContext(ctx, `
		SELECT id, beneficiary_code, beneficiary_name, account_number, ifsc_code, bank_name, email, mobile, payee_category
		FROM payees WHERE id = $1 AND is_deleted = FALSE`, id).
		Scan(&got.id, &got.beneficiaryCode, &got.beneficiaryName,
			&got.accNo, &got.ifsc, &got.bankName,
			&got.email, &got.mobile, &got.payeeCategory)
	if err == sql.ErrNoRows {
		got = nil
	}
	assert.Error(t, err, "expected error fetching soft deleted payee")
	assert.Nil(t, got, "soft deleted payee should not be returned by GetByID")
}
