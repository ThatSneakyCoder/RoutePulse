ALTER TABLE users
ADD COLUMN verify_email_token_hash TEXT,
ADD COLUMN verify_email_token_expires_at TIMESTAMP;

CREATE INDEX users_verify_email_token_idx
ON users (verify_email_token_expires_at)
WHERE verify_email_token_expires_at IS NOT NULL;
