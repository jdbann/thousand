-- name: GetVampire :one
SELECT
    *
FROM
    vampires
WHERE
    id = $1
LIMIT 1;

-- name: CreateVampire :one
INSERT INTO vampires (name)
    VALUES ($1)
RETURNING
    *;

-- name: GetVampires :many
SELECT
    *
FROM
    vampires;

