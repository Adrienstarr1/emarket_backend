-- +goose Up
CREATE TABLE users (
    index SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    age INT NOT NULL,
    email TEXT UNIQUE,
    password TEXT NOT NULL,
    id TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ

);

-- +goose down
DROP TABLE users;

