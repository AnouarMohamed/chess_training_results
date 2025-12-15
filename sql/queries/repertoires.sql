-- name: CreateRepertoire :one
INSERT INTO repertoires (player_id, name, color)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListRepertoiresByPlayer :many
SELECT * FROM repertoires
WHERE player_id = $1
ORDER BY created_at DESC;

-- name: GetRepertoireByID :one
SELECT * FROM repertoires
WHERE id = $1
LIMIT 1;

-- name: UpdateRepertoireName :one
UPDATE repertoires
SET name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteRepertoire :exec
DELETE FROM repertoires
WHERE id = $1;
