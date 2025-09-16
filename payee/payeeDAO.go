package payee

import (
	"context"
	"database/sql"
)

type PayeeRepository interface {
	Insert(context context.Context, p *payee) (int, error)
	GetByID(context context.Context, id int) (*payee, error)
}

type PayeeDB struct {
	db *sql.DB
}

func PostgresPayeeDB(db *sql.DB) *PayeeDB { //check function name and struct name
	return &PayeeDB{db: db}
}

func (r *PayeeDB) Insert(context context.Context, p *payee) (int, error) {
	query := `
		INSERT INTO payees (beneficiary_name, beneficiary_code, account_number,ifsc_code, bank_name, email, mobile, payee_category)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	var id int
	err := r.db.QueryRowContext(context, query,
		p.beneficiaryName,
		p.beneficiaryCode,
		p.accNo,
		p.ifsc,
		p.bankName,
		p.email,
		p.mobile,
		p.payeeCategory,
	).Scan(&id)
	return id, err
}

func (r *PayeeDB) GetByID(context context.Context, id int) (*payee, error) {
	query := `
		SELECT beneficiary_name, beneficiary_code, account_number,
		       ifsc_code, bank_name, email, mobile, payee_category
		FROM payees WHERE id=$1`
	row := r.db.QueryRowContext(context, query, id)

	var p payee
	err := row.Scan(
		&p.beneficiaryName, //repeating
		&p.beneficiaryCode,
		&p.accNo,
		&p.ifsc,
		&p.bankName,
		&p.email,
		&p.mobile,
		&p.payeeCategory,
	)
	p.id = id
	if err != nil {
		return nil, err
	}
	return &p, nil
}
