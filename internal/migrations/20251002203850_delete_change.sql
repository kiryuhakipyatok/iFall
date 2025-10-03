-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS iphones
DROP COLUMN IF EXISTS change
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS iphones
ADD COLUMN IF NOT EXISTS change NUMERIC NOT NULL DEFAULT 0
-- +goose StatementEnd
