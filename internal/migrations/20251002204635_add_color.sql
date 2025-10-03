-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS iphones
ADD COLUMN IF NOT EXISTS color VARCHAR(6) NOT NULL DEFAULT 'ffffff'
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS iphones
DROP COLUMN IF EXISTS color
-- +goose StatementEnd
