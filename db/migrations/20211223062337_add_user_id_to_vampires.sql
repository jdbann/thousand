-- +goose Up
-- +goose StatementBegin
ALTER TABLE vampires
    ADD COLUMN user_id UUID REFERENCES users (id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE vampires
    DROP COLUMN user_id;

-- +goose StatementEnd
