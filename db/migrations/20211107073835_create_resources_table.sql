-- +goose Up
-- +goose StatementBegin
CREATE TABLE resources (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    vampire_id uuid REFERENCES vampires (id) NOT NULL,
    description text NOT NULL,
    stationary boolean NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE resources;

-- +goose StatementEnd
