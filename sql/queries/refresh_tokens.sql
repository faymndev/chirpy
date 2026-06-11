-- name: CreateRefreshToken :one
insert into refresh_tokens (token, user_id, expires_at) 
values ($1, $2, $3)
returning *;


-- name: GetUserRefreshToken :one
select * from refresh_tokens
where user_id = $1
limit 1;

-- name: GetRefreshToken :one
select * from refresh_tokens
where token = $1
limit 1;

-- name: RevokeRefreshToken :exec 
update refresh_tokens
set revoked_at = now(), updated_at = now()
where token = $1;

-- name: SetRefreshToken :exec 
update refresh_tokens 
set token = $1, updated_at = now()
where user_id = $2;
