alter table users
add column updated_at TIMESTAMP NOT NULL DEFAULT now();