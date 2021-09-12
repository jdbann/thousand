-- +goose Up
-- +goose StatementBegin
CREATE TABLE vampires (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    name text NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE vampires;

-- +goose StatementEnd
