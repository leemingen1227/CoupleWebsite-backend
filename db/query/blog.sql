-- name: CreateBlog :one
INSERT INTO blog (id, pair_id, user_id, title, content, picture)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;