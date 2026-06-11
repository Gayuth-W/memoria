-- +goose Up
ALTER TABLE memories
ADD COLUMN importance_score FLOAT DEFAULT 0.5;

ALTER TABLE memories
ADD COLUMN embedding_hash TEXT;

-- +goose Down
ALTER TABLE memories DROP COLUMN importance_score;

ALTER TABLE memories DROP COLUMN embedding_hash;