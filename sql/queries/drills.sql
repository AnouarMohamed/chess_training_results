-- name: CreateDrillSession :one
INSERT INTO drill_sessions (player_id, target_line_id, status, current_ply)
VALUES ($1, $2, 'active', 0)
RETURNING *;

-- name: GetDrillSessionByID :one
SELECT * FROM drill_sessions
WHERE id = $1
LIMIT 1;

-- name: UpdateDrillProgress :one
UPDATE drill_sessions
SET current_ply = $2
WHERE id = $1
RETURNING *;

-- name: FinishDrillSession :one
UPDATE drill_sessions
SET status = 'finished',
    ended_at = now()
WHERE id = $1
RETURNING *;

-- name: AbortDrillSession :one
UPDATE drill_sessions
SET status = 'aborted',
    ended_at = now()
WHERE id = $1
RETURNING *;
