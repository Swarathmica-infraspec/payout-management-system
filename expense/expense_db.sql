CREATE TABLE expenses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    amount DECIMAL(14,2) NOT NULL,
    date_incurred DATE NOT NULL,
    category VARCHAR(50),
    notes TEXT,
    payee_id INT,
    status VARCHAR(20) NOT NULL DEFAULT 'Pending',
    receipt_uri VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

