-- +goose Up
-- +goose StatementBegin
CREATE TABLE iphones_new (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    price NUMERIC NOT NULL,
    change NUMERIC NOT NULL DEFAULT 0,
    color TEXT NOT NULL DEFAULT 'ffffff'
);

INSERT INTO iphones_new(id, name, price, change, color)
SELECT id, name, price, change, color FROM iphones;

DROP TABLE iphones;
ALTER TABLE iphones_new RENAME TO iphones;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE iphones_old (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    price NUMERIC NOT NULL,
    change NUMERIC NOT NULL DEFAULT 0,
    color TEXT NOT NULL DEFAULT 'ffffff'
);

INSERT INTO iphones_old(id, name, price, change, color)
SELECT id, name, price, change, color FROM iphones;

DROP TABLE iphones;
ALTER TABLE iphones_old RENAME TO iphones;
-- +goose StatementEnd
