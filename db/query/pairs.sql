-- name: CreatePair :one
INSERT INTO pairs (start_date)
VALUES ($1)
RETURNING *;