-- +goose Up
-- +goose StatementBegin
CREATE TYPE character_type AS enum (
    'mortal',
    'immortal'
);

CREATE TABLE characters (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    vampire_id uuid REFERENCES vampires (id) NOT NULL,
    name text NOT NULL,
    type character_type NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE characters;

DROP TABLE character_type;

-- +goose StatementEnd
