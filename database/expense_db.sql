CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    amount DECIMAL(14,2) NOT NULL,
    date_incurred DATE NOT NULL,
    category VARCHAR(50),
    notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'Pending',
    receipt_uri VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    payee_id INTEGER NOT NULL
);

