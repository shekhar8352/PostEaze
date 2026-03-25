-- Align users table with Firebase auth (fixes authenticate when 001 ran without firebase_id)
ALTER TABLE users ADD COLUMN IF NOT EXISTS firebase_id TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_firebase_id ON users(firebase_id);
