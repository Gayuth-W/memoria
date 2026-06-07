-- name: CreateUser :one
insert into users(id, api_key, created_at)
values ($1, $2, $3)
returning *;