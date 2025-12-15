-- 000001_init.down.sql

DROP TRIGGER IF EXISTS trg_lines_updated ON target_lines;
DROP TRIGGER IF EXISTS trg_openings_updated ON openings;
DROP TRIGGER IF EXISTS trg_repertoires_updated ON repertoires;
DROP TRIGGER IF EXISTS trg_players_updated ON players;

DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS drill_sessions;
DROP TABLE IF EXISTS target_lines;
DROP TABLE IF EXISTS openings;
DROP TABLE IF EXISTS repertoires;
DROP TABLE IF EXISTS players;

DROP TYPE IF EXISTS drill_status;
DROP TYPE IF EXISTS repertoire_color;

DROP EXTENSION IF EXISTS pgcrypto;
