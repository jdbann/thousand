-- name: CreateMemories :many
INSERT INTO memories (vampire_id)
SELECT
    unnest(@vampire_id::uuid[]) AS vampire_id
RETURNING
    *;

-- name: GetMemory :one
SELECT
    memories.*
FROM
    memories
    INNER JOIN vampires ON memories.vampire_id = vampires.id
WHERE
    vampires.id = @vampire_id
    AND memories.id = @memory_id;

-- name: GetMemoriesForVampire :many
SELECT
    *
FROM
    memories
WHERE
    memories.vampire_id = $1;

