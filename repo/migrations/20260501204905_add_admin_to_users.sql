-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN admin BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN admin;
-- +goose StatementEnd
