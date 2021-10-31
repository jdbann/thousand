-- name: CreateMemories :many
INSERT INTO memories (vampire_id)
SELECT
    unnest(@vampire_id::uuid[]) AS vampire_id
RETURNING
    *;

-- name: GetMemory :one
SELECT
    *
FROM
    memories
WHERE
    id = $1;

-- name: GetMemoriesForVampire :many
SELECT
    *
FROM
    memories
WHERE
    memories.vampire_id = $1;

