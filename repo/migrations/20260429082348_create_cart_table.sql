-- +goose Up
-- +goose StatementBegin
CREATE TABLE cart (
    index SERIAL PRIMARY KEY,
    user_id TEXT,
    product_id TEXT,
    quantity INT,
    Foreign Key (user_id) REFERENCES users(id),
    Foreign Key (product_id) REFERENCES products(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cart;
-- +goose StatementEnd
