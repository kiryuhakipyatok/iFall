-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS users
DROP COLUMN IF EXISTS chat_id;

ALTER TABLE IF EXISTS users
ADD COLUMN IF NOT EXISTS chat_id BIGINT UNIQUE
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS users
ADD COLUMN IF NOT EXISTS chat_id VARCHAR(16) UNIQUE
-- +goose StatementEnd
