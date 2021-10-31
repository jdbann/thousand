-- name: CreateExperience :one
INSERT INTO experiences (memory_id, description)
    VALUES ($1, $2)
RETURNING
    *;

-- name: GetExperiencesForVampire :many
SELECT
    experiences.*
FROM
    experiences
    INNER JOIN memories ON experiences.memory_id = memories.id
WHERE
    memories.vampire_id = $1;

