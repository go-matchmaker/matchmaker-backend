SET TIME ZONE 'Europe/Istanbul';

-- Create users table
CREATE TABLE users (
    id uuid PRIMARY KEY,
    role TEXT NOT NULL CHECK (role IN ('admin', 'user')) DEFAULT 'user',
    name            TEXT NOT NULL,
    surname         TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    phone_number    TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    created_at      timestamptz NOT NULL DEFAULT (now()),
    updated_at      timestamptz NOT NULL DEFAULT (now())
);
