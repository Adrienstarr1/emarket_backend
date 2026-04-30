-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
    index SERIAL PRIMARY KEY,
    id TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    quantity INT NOT NULL,
    user_id TEXT,
    Foreign Key (user_id) REFERENCES users(id),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
