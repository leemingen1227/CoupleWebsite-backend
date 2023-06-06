-- name: CreateUser :one
INSERT INTO users (id, email, password_digest, name, pair_id, is_email_verified)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET 
    email = COALESCE(sqlc.narg(email), email),
    password_digest = COALESCE(sqlc.narg(password_digest), password_digest),
    name = COALESCE(sqlc.narg(name), name),
    is_email_verified = COALESCE(sqlc.narg(is_email_verified), is_email_verified),
    pair_id = COALESCE(sqlc.narg(pair_id), pair_id),
    update_time = current_timestamp
WHERE id = sqlc.arg(id)
RETURNING *;
