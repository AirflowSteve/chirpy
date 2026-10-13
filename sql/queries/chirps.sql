-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
GROUP BY id
ORDER BY
created_at ASC;

-- name: GetChirpsFromAuthor :many
SELECT * FROM chirps WHERE user_id = $1
GROUP BY id
ORDER BY
created_at ASC;

-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE id = $1;

-- name: DeleteChirpByID :one
DELETE FROM chirps
WHERE id = $1
RETURNING *;