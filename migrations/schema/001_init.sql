-- +goose Up
create table users (
  id uuid primary key,
  api_key text unique not null,
  created_at timestamp default now()
);

create table sessions (
  id uuid primary key,
  user_id uuid references users(id) on delete cascade,
  title text,
  created_at timestamp default now()
);

create table memories (
  id uuid primary key,
  user_id uuid references users(id) on delete cascade,
  session_id uuid references sessions(id) on delete cascade,
  text text not null,
  created_at timestamp default now()
);

create index idx_memories_user on memories(user_id);
create index idx_memories_session on memories(session_id);

-- +goose Down
drop table if exists memories;
drop table if exists sessions;
drop table if exists users;