package payee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

type PayeeRepository interface {
	Insert(context context.Context, p *payee) (int, error)
	GetByID(context context.Context, id int) (*payee, error)
}

type payeeDB struct {
	db *sql.DB
}

func PayeeDB(db *sql.DB) *payeeDB {
	return &payeeDB{db: db}
}

var (
	ErrDuplicateCode    = errors.New("duplicate beneficiary code")
	ErrDuplicateAccount = errors.New("duplicate account number")
	ErrDuplicateEmail   = errors.New("duplicate email")
	ErrDuplicateMobile  = errors.New("duplicate mobile")
)

func (r *payeeDB) Insert(context context.Context, p *payee) (int, error) {
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
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			switch pgErr.Constraint {
			case "payees_beneficiary_code_key":
				return 0, ErrDuplicateCode
			case "payees_account_number_key":
				return 0, ErrDuplicateAccount
			case "payees_email_key":
				return 0, ErrDuplicateEmail
			case "payees_mobile_key":
				return 0, ErrDuplicateMobile
			}
		}
		return 0, fmt.Errorf("insert payee: %w", err)
	}
	return id, nil
}

func (r *payeeDB) GetByID(context context.Context, id int) (*payee, error) {
	query := `
		SELECT beneficiary_name, beneficiary_code, account_number,
		       ifsc_code, bank_name, email, mobile, payee_category
		FROM payees WHERE id=$1`
	row := r.db.QueryRowContext(context, query, id)

	var p payee
	err := row.Scan(
		&p.beneficiaryName,
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
