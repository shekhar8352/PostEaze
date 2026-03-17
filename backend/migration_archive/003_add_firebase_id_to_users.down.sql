-- Drop index first
DROP INDEX IF EXISTS idx_users_firebase_id;

-- Drop firebase_id column
ALTER TABLE users DROP COLUMN IF EXISTS firebase_id;
