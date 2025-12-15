-- name: CreatePlayer :one
INSERT INTO players (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPlayerByUsername :one
SELECT * FROM players
WHERE username = $1
LIMIT 1;

-- name: GetPlayerByEmail :one
SELECT * FROM players
WHERE email = $1
LIMIT 1;

-- name: UpdatePlayerLinks :one
UPDATE players
SET chesscom_link = $2,
    lichess_link  = $3,
    fide_id       = $4
WHERE id = $1
RETURNING *;

-- name: GetPlayerByID :one
SELECT * FROM players
WHERE id = $1
LIMIT 1;
