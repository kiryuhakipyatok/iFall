-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS iphones(
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(48) NOT NULL UNIQUE,
    price NUMERIC NOT NULL,
    change NUMERIC NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS iphones
-- +goose StatementEnd
