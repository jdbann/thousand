-- name: CreateSkill :one
INSERT INTO skills (vampire_id, description)
    VALUES (@vampire_id, @description)
RETURNING
    *;

-- name: GetSkillsForVampire :many
SELECT
    skills.*
FROM
    skills
WHERE
    skills.vampire_id = @vampire_id;

