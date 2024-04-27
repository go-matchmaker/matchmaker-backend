SET TIME ZONE 'Europe/Istanbul';

CREATE TYPE user_role AS ENUM (
    'admin',
    'customer'
);

-- Create users table
CREATE TABLE users (
    id uuid PRIMARY KEY,
    user_role       user_role NOT NULL,
    name            TEXT NOT NULL,
    surname         TEXT NOT NULL,
    email           TEXT NOT NULL,
    phone_number    TEXT NOT NULL,
    company_name    TEXT NOT NULL,
    company_type    TEXT NOT NULL,
    company_website TEXT NOT NULL,
    password_hash   TEXT NOT NULL,
    created_at      timestamptz NOT NULL DEFAULT (now()),
    updated_at      timestamptz NOT NULL DEFAULT (now())
);