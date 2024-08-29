-- +goose Up
CREATE TABLE IF NOT EXISTS avatar
(
    id        SERIAL PRIMARY KEY,
    user_uuid UUID NOT NULL,
    link      VARCHAR(255),
    create_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS avatar;
