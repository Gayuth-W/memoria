-- +goose Up
-- Add Full-Text Search support to memories

-- Step 1: Add tsvector column
ALTER TABLE memories ADD COLUMN IF NOT EXISTS tsv tsvector;

-- Step 2: Create GIN index for fast text search
CREATE INDEX IF NOT EXISTS idx_memories_fts ON memories USING GIN(tsv);

-- Step 3: Create trigger function
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION memories_tsv_trigger() RETURNS trigger AS $$
BEGIN
    new.tsv := to_tsvector('english', coalesce(new.text, ''));
    RETURN new;
END
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- Step 4: Create trigger
CREATE TRIGGER tsv_update
    BEFORE INSERT OR UPDATE ON memories
    FOR EACH ROW 
    EXECUTE FUNCTION memories_tsv_trigger();

-- Step 5: Backfill existing rows
UPDATE memories 
SET tsv = to_tsvector('english', coalesce(text, ''))
WHERE tsv IS NULL;

