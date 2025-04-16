-- name: CreatePost :one
INSERT INTO
    posts (user_id, title, content)
VALUES
    ($1, $2, $3)
RETURNING
    id,
    user_id,
    title,
    content,
    created_at,
    updated_at;

-- GetPostByID :one
SELECT
    id,
    user_id,
    title,
    content,
    created_at,
    updated_at
FROM
    posts
WHERE
    id = $1;

-- name: GetLatestPosts :many
SELECT
    p.id,
    p.title,
    p.content,
    p.created_at,
    u.username AS author
FROM
    posts p
    JOIN users u ON p.user_id = u.id
ORDER BY
    p.created_at DESC
LIMIT
    $2
OFFSET
    $1;