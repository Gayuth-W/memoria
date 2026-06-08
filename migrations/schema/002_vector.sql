-- +goose Up
-- Add Full-Text Search support to memories

-- Step 1: Add tsvector column
ALTER TABLE memories ADD COLUMN IF NOT EXISTS tsv tsvector;

-- Step 2: Create GIN index for fast text search
CREATE INDEX IF NOT EXISTS idx_memories_fts ON memories USING GIN(tsv);

-- Step 3: Create trigger function
CREATE OR REPLACE FUNCTION memories_tsv_trigger() RETURNS trigger AS $$
BEGIN
    new.tsv := to_tsvector('english', coalesce(new.text, ''));
    RETURN new;
END
$$ LANGUAGE plpgsql;

-- Step 4: Create trigger
DROP TRIGGER IF EXISTS tsv_update ON memories;
CREATE TRIGGER tsv_update
    BEFORE INSERT OR UPDATE ON memories
    FOR EACH ROW 
    EXECUTE FUNCTION memories_tsv_trigger();

-- Step 5: Backfill existing rows
UPDATE memories 
SET tsv = to_tsvector('english', coalesce(text, ''))
WHERE tsv IS NULL;


-- +goose Down
-- Remove Full-Text Search support only (not the tables!)

-- Step 1: Remove trigger
DROP TRIGGER IF EXISTS tsv_update ON memories;

-- Step 2: Remove trigger function
DROP FUNCTION IF EXISTS memories_tsv_trigger();

-- Step 3: Remove index
DROP INDEX IF EXISTS idx_memories_fts;

-- Step 4: Remove column
ALTER TABLE memories DROP COLUMN IF EXISTS tsv;