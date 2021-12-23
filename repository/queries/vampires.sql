-- name: GetVampire :one
SELECT
    *
FROM
    vampires
WHERE
    id = $1
LIMIT 1;

-- name: CreateVampire :one
INSERT INTO vampires (name, user_id)
    VALUES (@name, @user_id::uuid)
RETURNING
    *;

-- name: GetVampires :many
SELECT
    *
FROM
    vampires;

