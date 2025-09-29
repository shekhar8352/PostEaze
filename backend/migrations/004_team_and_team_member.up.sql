-- Migration: Add improvements to teams and team_members schema
-- ===============================================

-- ========== TEAMS TABLE UPDATES ==========

-- Add optional description and avatar/logo
ALTER TABLE teams
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS avatar_url TEXT;

-- Add visibility, status, settings, and owner role override
ALTER TABLE teams
    ADD COLUMN IF NOT EXISTS visibility TEXT CHECK (visibility IN ('private', 'public')) DEFAULT 'private',
    ADD COLUMN IF NOT EXISTS status TEXT CHECK (status IN ('active', 'archived', 'deleted')) DEFAULT 'active',
    ADD COLUMN IF NOT EXISTS settings JSONB DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS owner_role_override TEXT;

-- Indexes for teams
CREATE INDEX IF NOT EXISTS idx_teams_owner_id ON teams(owner_id);
-- Trigram index for name search (requires pg_trgm extension)
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS idx_teams_name_trgm ON teams USING gin (name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_active_teams ON teams(status) WHERE status = 'active';


-- ========== TEAM_MEMBERS TABLE UPDATES ==========

-- Add membership state, joined timestamp, invited_by, fine-grained permissions, primary team flag
ALTER TABLE team_members
    ADD COLUMN IF NOT EXISTS status TEXT CHECK (status IN ('active','invited','pending','removed')) DEFAULT 'active',
    ADD COLUMN IF NOT EXISTS joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN IF NOT EXISTS invited_by UUID REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS permissions JSONB DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS is_primary BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS last_active_at TIMESTAMP;

-- Indexes for team_members
CREATE INDEX IF NOT EXISTS idx_team_members_team_id ON team_members(team_id);
CREATE INDEX IF NOT EXISTS idx_team_members_user_id ON team_members(user_id);
CREATE INDEX IF NOT EXISTS idx_team_members_team_user ON team_members(team_id, user_id);
CREATE INDEX IF NOT EXISTS idx_active_team_members ON team_members(team_id) WHERE status = 'active';
