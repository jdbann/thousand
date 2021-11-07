-- name: CreateMark :one
INSERT INTO marks (vampire_id, description)
    VALUES (@vampire_id, @description)
RETURNING
    *;

-- name: GetMarksForVampire :many
SELECT
    marks.*
FROM
    marks
WHERE
    marks.vampire_id = @vampire_id;

