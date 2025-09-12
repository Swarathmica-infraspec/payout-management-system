package expense

import (
	"context"
	"database/sql"
	"log"
)

type ExpenseRepository interface {
	Insert(ctx context.Context, e *expense) (int, error)
	GetByID(ctx context.Context, id int) (*expense, error)
	List(context context.Context) ([]expense, error)
}

type ExpensePostgresDB struct {
	Db *sql.DB
}

func NewPostgresExpenseDB(db *sql.DB) *ExpensePostgresDB {
	return &ExpensePostgresDB{Db: db}
}

func (r *ExpensePostgresDB) Insert(ctx context.Context, e *expense) (int, error) {
	query := `
		INSERT INTO expenses 
		(title, amount, date_incurred, category, notes, payee_id, receipt_uri,status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id`
	var id int
	err := r.Db.QueryRowContext(ctx, query,
		e.title,
		e.amount,
		e.dateIncurred,
		e.category,
		e.notes,
		e.payeeID,
		e.receiptURI,
		e.status,
	).Scan(&id)
	return id, err
}

func (r *ExpensePostgresDB) GetByID(ctx context.Context, id int) (*expense, error) {
	query := `
		SELECT title, amount, date_incurred, category, notes, payee_id, receipt_uri 
		FROM expenses WHERE id=$1`
	row := r.Db.QueryRowContext(ctx, query, id)
	var e expense
	err := row.Scan(
		&e.title,
		&e.amount,
		&e.dateIncurred,
		&e.category,
		&e.notes,
		&e.payeeID,
		&e.receiptURI,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *ExpensePostgresDB) List(context context.Context) ([]expense, error) {
	rows, err := s.Db.QueryContext(context, `
        SELECT title, amount, date_incurred, category, notes, payee_id, receipt_uri 
		FROM expenses
        ORDER BY payee_id ASC
    `)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var expenses []expense
	for rows.Next() {
		var e expense
		err := rows.Scan(&e.title,
			&e.amount,
			&e.dateIncurred,
			&e.category,
			&e.notes,
			&e.payeeID,
			&e.receiptURI)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}

	return expenses, nil
}

type ExpenseWithPayee struct {
	ExpenseID       int
	Title           string
	Amount          float64
	DateIncurred    string
	Notes           string
	BeneficiaryName string
	BeneficiaryCode string
	AccountNumber   int
	IFSCCode        string
	BankName        string
	Email           string
}

func (r *ExpensePostgresDB) ListExpensesForPayout(ctx context.Context) ([]ExpenseWithPayee, float64, error) {
	query := `
        SELECT e.id, e.title, e.amount, e.date_incurred, e.notes, p.beneficiary_name, p.beneficiary_code, p.account_number, p.ifsc_code, p.bank_name, p.email
        FROM expenses e
        JOIN payees p ON e.payee_id = p.id
        WHERE e.status= 'Pending'
        ORDER BY e.date_incurred DESC;
    `
	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var list []ExpenseWithPayee
	var total float64
	for rows.Next() {
		var ew ExpenseWithPayee
		if err := rows.Scan(
			&ew.ExpenseID,
			&ew.Title,
			&ew.Amount,
			&ew.DateIncurred,
			&ew.Notes,
			&ew.BeneficiaryName,
			&ew.BeneficiaryCode,
			&ew.AccountNumber,
			&ew.IFSCCode,
			&ew.BankName,
			&ew.Email,
		); err != nil {
			return nil, 0, err
		}
		total += ew.Amount
		list = append(list, ew)
	}
	return list, total, nil
}
