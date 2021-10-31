-- name: CreateExperience :one
INSERT INTO experiences (memory_id, description)
    VALUES ((
            SELECT
                memories.id
            FROM
                memories
            WHERE
                memories.id = @memory_id
                AND memories.vampire_id = @vampire_id),
            @description)
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

