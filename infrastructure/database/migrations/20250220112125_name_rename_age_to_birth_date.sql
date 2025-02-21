-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_metrics RENAME COLUMN age TO birth_date;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_metrics RENAME COLUMN birth_date TO age;
-- +goose StatementEnd
