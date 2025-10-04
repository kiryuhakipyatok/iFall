-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS iphones (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    price NUMERIC NOT NULL,
    change NUMERIC NOT NULL DEFAULT 0,
    color TEXT NOT NULL DEFAULT 'ffffff'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS iphones;
-- +goose StatementEnd
