package payoutmanagementsystem

import (
	"context"
	"database/sql"
)

type ExpenseRepository interface {
	Create(ctx context.Context, e *expense) (int, error)
	GetByID(ctx context.Context, id int) (*expense, error)
}

type ExpensePostgresDB struct {
	db *sql.DB
}

func NewPostgresExpenseDB(db *sql.DB) *ExpensePostgresDB {
	return &ExpensePostgresDB{db: db}
}

func (r *ExpensePostgresDB) Insert(ctx context.Context, e *expense) (int, error) {
	query := `
		INSERT INTO expenses
		(title, amount, date_incurred, category, notes, payee_id, receipt_uri)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, query,
		e.title,
		e.amount,
		e.dateIncurred,
		e.category,
		e.notes,
		e.payeeID,
		e.receiptURI,
	).Scan(&id)
	return id, err
}

func (r *ExpensePostgresDB) GetByID(ctx context.Context, id int) (*expense, error) {
	query := `
		SELECT title, amount, date_incurred, category, notes, payee_id, receipt_uri
		FROM expenses WHERE id=$1`
	row := r.db.QueryRowContext(ctx, query, id)
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
