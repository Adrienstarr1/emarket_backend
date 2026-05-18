-- +goose Up
-- +goose StatementBegin
ALTER TABLE products ADD COLUMN price INT;
ALTER TABLE products ADD COLUMN description TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products DROP COLUMN price;
ALTER TABLE products DROP COLUMN description;
-- +goose StatementEnd
