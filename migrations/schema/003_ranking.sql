-- +goose Up
ALTER TABLE memories
ADD COLUMN importance_score FLOAT DEFAULT 0.5;

ALTER TABLE memories
ADD COLUMN embedding_hash TEXT;
