-- name: CreateUser :one
insert into users (id, email)
values (gen_random_uuid(), $1)
returning *;

-- name: Reset :exec
delete from users;
