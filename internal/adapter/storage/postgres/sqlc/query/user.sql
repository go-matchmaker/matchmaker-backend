-- name: InsertUser :one
INSERT INTO users (
    id,
    user_role,
    name,
    surname,
    email,
    phone_number,
    company_name,
    company_type,
    company_website,
    password_hash,
    created_at,
    updated_at
 ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
    name = COALESCE(sqlc.narg(name), name),
    surname = COALESCE(sqlc.narg(surname), surname),
    email = COALESCE(sqlc.narg(email), email),
    phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
    company_name = COALESCE(sqlc.narg(company_name), company_name),
    company_type = COALESCE(sqlc.narg(company_type), company_type),
    company_website = COALESCE(sqlc.narg(company_website), company_website),
    password_hash = COALESCE(sqlc.narg(password_hash), password_hash),
    updated_at = COALESCE(sqlc.narg(updated_at), updated_at)
WHERE
id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;