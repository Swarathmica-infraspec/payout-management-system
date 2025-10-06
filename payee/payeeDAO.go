package payee

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PayeeRepository interface {
	Insert(ctx context.Context, p *payee) (int, error)
	GetByID(ctx context.Context, id int) (*payee, error)
	List(ctx context.Context) ([]payee, error)
	Update(ctx context.Context, p *payee) (*payee, error)
	
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

func (r *payeeDB) Insert(ctx context.Context, p *payee) (int, error) {
	query := `
        INSERT INTO payees (beneficiary_name, beneficiary_code, account_number,ifsc_code, bank_name, email, mobile, payee_category)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, query,
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
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.ConstraintName {
			case "payees_beneficiary_code_key":
				return 0, ErrDuplicateCode
			case "payees_account_number_key":
				return 0, ErrDuplicateAccount
			case "payees_email_key":
				return 0, ErrDuplicateEmail
			case "payees_mobile_key":
				return 0, ErrDuplicateMobile
			}
			return 0, fmt.Errorf("insert payee: %w", err)
		}
		return 0, err
	}
	return id, nil
}

func (r *payeeDB) GetByID(ctx context.Context, id int) (*payee, error) {
	query := `
        SELECT beneficiary_name, beneficiary_code, account_number,
               ifsc_code, bank_name, email, mobile, payee_category
        FROM payees WHERE id=$1`
	row := r.db.QueryRowContext(ctx, query, id)

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
		return nil, fmt.Errorf("get payee by id %d: %w", id, err)
	}
	return &p, nil
}
func (r *payeeDB) List(ctx context.Context) ([]payee, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT id, beneficiary_name, beneficiary_code, account_number, ifsc_code, bank_name, email, mobile, payee_category
        FROM payees
        ORDER BY id ASC
    `)
	if err != nil {
		return nil, fmt.Errorf("List payee: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var payees []payee
	for rows.Next() {
		var p payee
		err := rows.Scan(&p.id, &p.beneficiaryName, &p.beneficiaryCode, &p.accNo, &p.ifsc,
			&p.bankName,
			&p.email,
			&p.mobile,
			&p.payeeCategory)
		if err != nil {
			return nil, err
		}
		payees = append(payees, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payees, nil
}

func (s *payeeDB) Update(ctx context.Context, p *payee) (*payee, error) {
	query := `
        UPDATE payees SET
            beneficiary_name = $1,
            beneficiary_code = $2,
            account_number   = $3,
            ifsc_code        = $4,
            bank_name        = $5,
            email            = $6,
            mobile           = $7,
            payee_category   = $8
        WHERE id = $9
        RETURNING id, beneficiary_name, beneficiary_code, account_number, ifsc_code, bank_name, email, mobile, payee_category
    `

	var updatedPayee payee
	err := s.db.QueryRowContext(ctx, query,
		p.beneficiaryName,
		p.beneficiaryCode,
		p.accNo,
		p.ifsc,
		p.bankName,
		p.email,
		p.mobile,
		p.payeeCategory,
		p.id,
	).Scan(
		&updatedPayee.id,
		&updatedPayee.beneficiaryName,
		&updatedPayee.beneficiaryCode,
		&updatedPayee.accNo,
		&updatedPayee.ifsc,
		&updatedPayee.bankName,
		&updatedPayee.email,
		&updatedPayee.mobile,
		&updatedPayee.payeeCategory,
	)
	if err != nil {
		return nil, err
	}

	return &updatedPayee, nil
}
