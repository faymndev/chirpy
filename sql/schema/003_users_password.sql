-- +goose Up
alter table users 
  add column password text not null default '';

-- +goose Down
alter table users
  drop column password;
