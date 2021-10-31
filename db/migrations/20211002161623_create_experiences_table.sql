-- +goose Up
-- +goose StatementBegin
CREATE TABLE experiences (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    memory_id uuid REFERENCES memories (id) NOT NULL,
    description text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE experiences;

-- +goose StatementEnd
