-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS temp_users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(50) NOT NULL,
    hashed_password VARCHAR(255) NOT NULL,
    token VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS temp_users;
-- +goose StatementEnd
