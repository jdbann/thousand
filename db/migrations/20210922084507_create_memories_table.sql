-- +goose Up
-- +goose StatementBegin
CREATE TABLE memories (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    vampire_id uuid REFERENCES vampires (id) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE memories;

-- +goose StatementEnd
