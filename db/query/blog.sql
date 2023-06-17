-- name: CreateBlog :one
INSERT INTO blog (id, pair_id, user_id, title, content, picture)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetBlogByBlogID :one
SELECT * FROM blog
WHERE id = $1 LIMIT 1;

-- name: CountBlogsByPairID :one
SELECT COUNT(*) FROM blog
WHERE pair_id = $1;

-- name: GetBlogsByPairID :many
SELECT * FROM blog
WHERE pair_id = $1 ORDER BY create_time DESC LIMIT $2 OFFSET $3;