-- name: CreateGame :one
INSERT INTO games (player_id, opening_id, target_line_id, drill_session_id, pgn_raw, result, played_at)
VALUES ($1, $2, $3, $4, $5, $6, COALESCE($7, now()))
RETURNING *;

-- name: GetGameByID :one
SELECT * FROM games
WHERE id = $1
LIMIT 1;

-- name: ListGamesByPlayerLimit :many
SELECT * FROM games
WHERE player_id = $1
ORDER BY played_at DESC
LIMIT $2;

-- name: ListGamesByPlayerAndOpeningLimit :many
SELECT * FROM games
WHERE player_id = $1
  AND opening_id = $2
ORDER BY played_at DESC
LIMIT $3;
