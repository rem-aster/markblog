-- name: CreateComment :one
INSERT INTO
    comments (post_id, user_id, content)
VALUES
    ($1, $2, $3)
RETURNING
    id,
    post_id,
    user_id,
    content,
    created_at,
    updated_at;

-- name: GetCommentByID :one
SELECT
    id,
    post_id,
    user_id,
    content,
    created_at,
    updated_at
FROM
    comments
WHERE
    id = $1;

-- name: GetLatestCommentsForPost :many
SELECT
    c.id,
    c.content,
    c.created_at,
    u.username AS username
FROM
    comments c
JOIN
    users u ON c.user_id = u.id
WHERE
    c.post_id = $1
ORDER BY
    c.created_at DESC
LIMIT
    $2
OFFSET
    $3;