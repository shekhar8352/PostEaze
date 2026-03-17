-- Rollback Migration: Remove improvements from teams and team_members
-- ===============================================

-- ========== TEAMS TABLE ROLLBACK ==========

ALTER TABLE teams
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS avatar_url,
    DROP COLUMN IF EXISTS visibility,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS settings,
    DROP COLUMN IF EXISTS owner_role_override;

-- Drop indexes for teams
DROP INDEX IF EXISTS idx_teams_owner_id;
DROP INDEX IF EXISTS idx_teams_name_trgm;
DROP INDEX IF EXISTS idx_active_teams;


-- ========== TEAM_MEMBERS TABLE ROLLBACK ==========

ALTER TABLE team_members
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS joined_at,
    DROP COLUMN IF EXISTS invited_by,
    DROP COLUMN IF EXISTS permissions,
    DROP COLUMN IF EXISTS is_primary,
    DROP COLUMN IF EXISTS last_active_at;

-- Drop indexes for team_members
DROP INDEX IF EXISTS idx_team_members_team_id;
DROP INDEX IF EXISTS idx_team_members_user_id;
DROP INDEX IF EXISTS idx_team_members_team_user;
DROP INDEX IF EXISTS idx_active_team_members;
