-- name: GetMemory :one
SELECT
    *
FROM
    memories
WHERE
    id = $1
LIMIT 1;

-- name: GetMemoriesForVampire :many
SELECT
    *
FROM
    memories
WHERE
    vampire_id = $1;

-- name: CreateMemory :one
INSERT INTO memories (vampire_id)
    VALUES ($1)
RETURNING
    *;

-- name: CreateMemories :many
INSERT INTO memories (vampire_id)
SELECT
    unnest(@vampire_id::uuid[]) AS vampire_id
RETURNING
    *;

