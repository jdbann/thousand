-- name: CreateUser :one
INSERT INTO users (email, password_hash)
    VALUES (lower(@email), crypt(@password::text, gen_salt('bf', 8)))
RETURNING
    id, email;

