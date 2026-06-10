-- name: CreateChirp :one
insert into chirps (id, user_id, body)
values (gen_random_uuid(), $1, $2)
returning *;
