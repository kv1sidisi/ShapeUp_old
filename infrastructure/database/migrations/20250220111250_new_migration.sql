-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_sessions DROP COLUMN last_activity;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_sessions ADD COLUMN last_activity TIMESTAMP;
-- +goose StatementEnd
