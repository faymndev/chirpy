-- name: CreateRefreshToken :one
insert into refresh_tokens (token, user_id, expires_at) 
values ($1, $2, $3)
returning *;


-- name: GetRefreshToken :one
select * from refresh_tokens
where id = $1
limit 1;

-- name: RevokeRefreshToken :exec 
update refresh_tokens
set revoked_at = $1
where id = $2;
