-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, hashed_password)
values (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
returning *;

-- name: GetUserById :one
select * from users
where id = $1;

-- name: GetUserByEmail :one
select * from users
where email = $1;
