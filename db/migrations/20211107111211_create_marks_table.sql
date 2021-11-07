-- +goose Up
-- +goose StatementBegin
CREATE TABLE marks (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    vampire_id uuid REFERENCES vampires (id) NOT NULL,
    description text NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE marks;

-- +goose StatementEnd
