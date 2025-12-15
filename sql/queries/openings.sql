-- name: CreateOpening :one
INSERT INTO openings (repertoire_id, name, eco)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListOpeningsByRepertoire :many
SELECT * FROM openings
WHERE repertoire_id = $1
ORDER BY created_at DESC;

-- name: GetOpeningByID :one
SELECT * FROM openings
WHERE id = $1
LIMIT 1;

-- name: UpdateOpening :one
UPDATE openings
SET name = COALESCE($2, name),
    eco  = COALESCE($3, eco)
WHERE id = $1
RETURNING *;

-- name: DeleteOpening :exec
DELETE FROM openings
WHERE id = $1;
