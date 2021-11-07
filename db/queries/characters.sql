-- name: CreateCharacter :one
INSERT INTO characters (vampire_id, name, type)
    VALUES (@vampire_id, @name, @type)
RETURNING
    *;

-- name: GetCharactersForVampire :many
SELECT
    characters.*
FROM
    characters
WHERE
    characters.vampire_id = @vampire_id;

