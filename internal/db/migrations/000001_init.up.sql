-- 000001_init.up.sql

-- UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Enums
DO $$ BEGIN
  CREATE TYPE repertoire_color AS ENUM ('white', 'black');
EXCEPTION
  WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
  CREATE TYPE drill_status AS ENUM ('active', 'finished', 'aborted');
EXCEPTION
  WHEN duplicate_object THEN null;
END $$;

-- Players
CREATE TABLE IF NOT EXISTS players (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username        TEXT NOT NULL UNIQUE,
  email           TEXT UNIQUE,
  password_hash   TEXT NOT NULL,
  chesscom_link   TEXT,
  lichess_link    TEXT,
  fide_id         TEXT,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Repertoires
CREATE TABLE IF NOT EXISTS repertoires (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id   UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  name        TEXT NOT NULL,
  color       repertoire_color NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (player_id, name, color)
);

CREATE INDEX IF NOT EXISTS idx_repertoires_player ON repertoires(player_id);

-- Openings
CREATE TABLE IF NOT EXISTS openings (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  repertoire_id UUID NOT NULL REFERENCES repertoires(id) ON DELETE CASCADE,
  name          TEXT NOT NULL,
  eco           TEXT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (repertoire_id, name)
);

CREATE INDEX IF NOT EXISTS idx_openings_repertoire ON openings(repertoire_id);

-- Target lines (strict SAN sequence)
CREATE TABLE IF NOT EXISTS target_lines (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  opening_id  UUID NOT NULL REFERENCES openings(id) ON DELETE CASCADE,
  name        TEXT NOT NULL,
  moves_san   JSONB NOT NULL,          -- array of SAN strings, in order
  start_fen   TEXT,                    -- nullable, default is initial position
  is_primary  BOOLEAN NOT NULL DEFAULT FALSE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (opening_id, name)
);

CREATE INDEX IF NOT EXISTS idx_lines_opening ON target_lines(opening_id);

-- Drill sessions
CREATE TABLE IF NOT EXISTS drill_sessions (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id     UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  target_line_id UUID NOT NULL REFERENCES target_lines(id) ON DELETE CASCADE,
  status        drill_status NOT NULL DEFAULT 'active',
  current_ply   INT NOT NULL DEFAULT 0,
  started_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  ended_at      TIMESTAMPTZ,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_drills_player ON drill_sessions(player_id);
CREATE INDEX IF NOT EXISTS idx_drills_line ON drill_sessions(target_line_id);

-- Games
CREATE TABLE IF NOT EXISTS games (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_id       UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
  opening_id      UUID REFERENCES openings(id) ON DELETE SET NULL,
  target_line_id  UUID REFERENCES target_lines(id) ON DELETE SET NULL,
  drill_session_id UUID REFERENCES drill_sessions(id) ON DELETE SET NULL,
  pgn_raw         TEXT NOT NULL,
  result          TEXT NOT NULL DEFAULT '*', -- '1-0'|'0-1'|'1/2-1/2'|'*'
  played_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_games_player ON games(player_id);
CREATE INDEX IF NOT EXISTS idx_games_opening ON games(opening_id);
CREATE INDEX IF NOT EXISTS idx_games_played_at ON games(played_at DESC);

-- Reviews (1–1 with game)
CREATE TABLE IF NOT EXISTS reviews (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  game_id       UUID NOT NULL UNIQUE REFERENCES games(id) ON DELETE CASCADE,
  analysis_pgn  TEXT NOT NULL,
  summary_text  TEXT NOT NULL,
  deviations    JSONB NOT NULL DEFAULT '[]'::jsonb,  -- [{ply, expected, played}]
  eval_events   JSONB NOT NULL DEFAULT '[]'::jsonb,  -- [{ply, eval_before, eval_after, delta}]
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Updated_at trigger helper
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach updated_at triggers
DO $$ BEGIN
  CREATE TRIGGER trg_players_updated
  BEFORE UPDATE ON players
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
EXCEPTION WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
  CREATE TRIGGER trg_repertoires_updated
  BEFORE UPDATE ON repertoires
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
EXCEPTION WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
  CREATE TRIGGER trg_openings_updated
  BEFORE UPDATE ON openings
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
EXCEPTION WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
  CREATE TRIGGER trg_lines_updated
  BEFORE UPDATE ON target_lines
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
EXCEPTION WHEN duplicate_object THEN null;
END $$;
