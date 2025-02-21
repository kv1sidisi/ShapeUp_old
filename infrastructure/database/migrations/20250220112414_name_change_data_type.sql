-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_metrics
    ALTER COLUMN birth_date TYPE DATE
        USING TO_DATE(birth_date::TEXT, 'YYYYMMDD');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_metrics
    ALTER COLUMN birth_date TYPE INTEGER
        USING EXTRACT(YEAR FROM birth_date) * 10000 + EXTRACT(MONTH FROM birth_date) * 100 + EXTRACT(DAY FROM birth_date);
-- +goose StatementEnd