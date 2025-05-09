-- name: CreatePost :one
INSERT INTO
    posts (user_id, content)
VALUES
    ($1, $2)
RETURNING
    id,
    user_id,
    content,
    created_at,
    updated_at;

-- name: GetPostByID :one
SELECT
    id,
    user_id,
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
    p.content,
    p.created_at,
    u.username
FROM 
    posts p
JOIN 
    users u ON p.user_id = u.id
ORDER BY 
    p.created_at DESC
LIMIT 
    $1
OFFSET 
    $2;