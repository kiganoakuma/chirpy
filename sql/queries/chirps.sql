-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
  gen_random_uuid(),
  NOW(),
  NOW(),
  $1,
  $2
)
returning *;

-- name: GetAllChirps :many
select * from chirps
order by created_at asc;

-- name: GetChirpById :one
select * from chirps
where id = $1;


-- name: DelChirpById :exec
delete from chirps
where id = $1 and user_id = $2;
