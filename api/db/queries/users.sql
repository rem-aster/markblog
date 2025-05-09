-- name: CreateUser :one
INSERT INTO
    users (username, password_hash)
VALUES
    ($1, $2)
RETURNING
    id,
    username,
    password_hash,
    created_at,
    updated_at;

-- name: GetUserByID :one
SELECT
    id,
    username,
    password_hash,
    created_at,
    updated_at
FROM
    users
WHERE
    id = $1;

-- name: GetUserByUsername :one
SELECT
    id,
    username,
    password_hash,
    created_at,
    updated_at
FROM
    users
WHERE
    username = $1;

-- name: CheckUserExists :one
SELECT
    EXISTS (
        SELECT
            1
        FROM
            users
        WHERE
            username = $1
    ) AS user_exists;

-- name: GetLatestUserActivity :many
WITH user_info AS (
    SELECT id
    FROM users
    WHERE username = $1
)
SELECT
    p.id AS post_id,
    p.content,
    p.created_at AS action_time,
    'post' AS action_type
FROM
    posts p
WHERE
    p.user_id = (SELECT id FROM user_info)
UNION ALL
SELECT
    c.id AS post_id,
    c.content,
    c.created_at AS action_time,
    'comment' AS action_type
FROM
    comments c
WHERE
    c.user_id = (SELECT id FROM user_info)
ORDER BY
    action_time DESC
LIMIT 
    $2
OFFSET 
    $3;