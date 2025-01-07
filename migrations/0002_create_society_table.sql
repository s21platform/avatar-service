-- +goose Up
CREATE TABLE IF NOT EXISTS society
(
    id        SERIAL PRIMARY KEY,
    uuid      UUID NOT NULL,
    link      VARCHAR(255),
    create_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS society;
