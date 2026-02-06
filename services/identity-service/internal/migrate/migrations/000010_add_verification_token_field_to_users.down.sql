DROP INDEX IF EXISTS users_verify_email_token_idx;

ALTER TABLE users
DROP COLUMN IF EXISTS verify_email_token_hash,
DROP COLUMN IF EXISTS verify_email_token_expires_at;
