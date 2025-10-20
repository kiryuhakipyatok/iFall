-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN desired_price NUMERIC NOT NULL DEFAULT 0
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN desired_price
-- +goose StatementEnd
