-- name: Insert :one
INSERT INTO users (
    id,
    role,
    name,
    surname,
    email,
    phone_number,
    password_hash,
    created_at,
    updated_at
 ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: Update :one
UPDATE users
SET
    name = COALESCE(sqlc.narg(name), name),
    surname = COALESCE(sqlc.narg(surname), surname),
    email = COALESCE(sqlc.narg(email), email),
    phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
    password_hash = COALESCE(sqlc.narg(password_hash), password_hash),
    updated_at = COALESCE(sqlc.narg(updated_at), updated_at)
WHERE
id = sqlc.arg(id)
RETURNING *;

-- name: DeleteOne :exec
DELETE FROM users
WHERE id = $1;

-- name: DeleteAll :exec
DELETE FROM users;