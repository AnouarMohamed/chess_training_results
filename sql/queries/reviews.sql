-- name: CreateReview :one
INSERT INTO reviews (game_id, analysis_pgn, summary_text, deviations, eval_events)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetReviewByGameID :one
SELECT * FROM reviews
WHERE game_id = $1
LIMIT 1;

-- name: GetReviewByID :one
SELECT * FROM reviews
WHERE id = $1
LIMIT 1;
