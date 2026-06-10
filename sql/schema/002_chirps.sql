-- +goose Up
create table chirps (
  id uuid primary key,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now(),
  user_id uuid not null references users (id) on delete cascade,
  body text not null
);

-- +goose Down 
drop table chirps;
