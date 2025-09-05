package payee

import (
	"context"
	"database/sql"
	"log"
)

type PayeeRepository interface {
	Insert(context context.Context, p *payee) (int, error)
	GetByID(context context.Context, id int) (*payee, error)
	List(ctx context.Context) ([]payee, error)
	Update(ctx context.Context) (*payee, error)
	Delete(ctx context.Context, id int) error
}

type PayeePostgresDB struct {
	db *sql.DB
}

func PostgresPayeeDB(db *sql.DB) *PayeePostgresDB {
	return &PayeePostgresDB{db: db}
}

func (r *PayeePostgresDB) Insert(context context.Context, p *payee) (int, error) {
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

func (r *PayeePostgresDB) GetByID(context context.Context, id int) (*payee, error) {
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

func (s *PayeePostgresDB) List(context context.Context) ([]payee, error) {
	rows, err := s.db.QueryContext(context, `
        SELECT id, beneficiary_name, beneficiary_code, account_number, ifsc_code, bank_name, email, mobile, payee_category
        FROM payees
        ORDER BY id ASC
    `)
	if err != nil {
		return nil, err
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

	return payees, nil
}

func (s *PayeePostgresDB) Update(ctx context.Context, p *payee) (*payee, error) {
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

func (r *PayeePostgresDB) Delete(context context.Context, id int) error {
	_, err := r.db.ExecContext(context, "DELETE FROM payees WHERE id=$1", id)
	return err
}
