-- name: CreateUsersPair :one
INSERT INTO user_pairs (pair_id, user_id)
VALUES ($1, $2)
RETURNING *;
