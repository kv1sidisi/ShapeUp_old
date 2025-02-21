-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    is_confirmed BOOLEAN DEFAULT FALSE
);

CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uid UUID NOT NULL,
    refresh_token TEXT NOT NULL,
    device_info TEXT,
    ip_address TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL,
    last_activity TIMESTAMP,
    CONSTRAINT fk_user_sessions_uid FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE user_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uid UUID NOT NULL,
    name TEXT,
    height NUMERIC(5,1),
    weight NUMERIC(5,1),
    age INTEGER,
    gender TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user_metrics_uid FOREIGN KEY (uid) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS user_metrics;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
