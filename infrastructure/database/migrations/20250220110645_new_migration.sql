-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_sessions DROP COLUMN expires_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_sessions ADD COLUMN expires_at TIMESTAMP;
-- +goose StatementEnd
