-- name: CreateChirp :one
insert into chirps (id, user_id, body)
values (gen_random_uuid(), $1, $2)
returning *;

-- name: GetChirps :many
select * from chirps
order by created_at asc;

-- name: GetChirp :one
select * from chirps 
where id = $1
limit 1;

-- name: DeleteChirp :exec
delete from chirps
where id = $1;
