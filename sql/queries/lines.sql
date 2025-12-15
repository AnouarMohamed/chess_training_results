-- name: CreateTargetLine :one
INSERT INTO target_lines (opening_id, name, moves_san, start_fen, is_primary)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListTargetLinesByOpening :many
SELECT * FROM target_lines
WHERE opening_id = $1
ORDER BY is_primary DESC, created_at DESC;

-- name: GetTargetLineByID :one
SELECT * FROM target_lines
WHERE id = $1
LIMIT 1;

-- name: UpdateTargetLine :one
UPDATE target_lines
SET name      = COALESCE($2, name),
    moves_san = COALESCE($3, moves_san),
    start_fen = COALESCE($4, start_fen),
    is_primary= COALESCE($5, is_primary)
WHERE id = $1
RETURNING *;

-- name: DeleteTargetLine :exec
DELETE FROM target_lines
WHERE id = $1;
