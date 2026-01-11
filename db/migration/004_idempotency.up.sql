CREATE TABLE idempotency_keys (
    key VARCHAR(100) PRIMARY KEY,
    response_status INT NOT NULL,
    response_body JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
