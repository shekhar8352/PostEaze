DROP INDEX IF EXISTS idx_users_firebase_id;
ALTER TABLE users DROP COLUMN IF EXISTS firebase_id;
