-- name: CreateUser :one
insert into users (id, email, password)
values (gen_random_uuid(), $1, $2)
returning *;

-- name: GetUser :one
select * from users 
where email = $1
limit 1;

-- name: GetUserById :one
select * from users 
where id = $1
limit 1;

-- name: Reset :exec
delete from users;

-- name: UpdateUser :one
update users  
set email = $1, password = $2, updated_at = now()
where id = $3
returning *;

-- name: UpgradeUser :one
update users  
set is_chirpy_red = $1, updated_at = now()
where id = $2
returning *;
