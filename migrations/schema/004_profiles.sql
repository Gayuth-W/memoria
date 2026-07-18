-- +goose Up
CREATE TABLE pinned_facts (
    id uuid primary key,
    user_id uuid references users(id) on delete cascade,
    fact text not null,
    created_at timestamp default now(),
    UNIQUE(user_id, fact)
);
