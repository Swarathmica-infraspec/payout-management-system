CREATE TABLE IF NOT EXISTS payees (
    id  INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    beneficiary_name TEXT NOT NULL,
    beneficiary_code TEXT NOT NULL UNIQUE,
    account_number BIGINT NOT NULL UNIQUE,
    ifsc_code VARCHAR(11) NOT NULL,
    bank_name VARCHAR(50) NOT NULL,
    email TEXT NOT NULL UNIQUE,
    mobile BIGINT NOT NULL UNIQUE,
    payee_category TEXT
);

ALTER TABLE payees ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;
