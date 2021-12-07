-- name: CreateResource :one
INSERT INTO resources (vampire_id, description, stationary)
    VALUES (@vampire_id, @description, @stationary)
RETURNING
    *;

-- name: GetResourcesForVampire :many
SELECT
    resources.*
FROM
    resources
WHERE
    resources.vampire_id = @vampire_id;

