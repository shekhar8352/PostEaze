-- Add firebase_id column to users table
ALTER TABLE users ADD COLUMN firebase_id TEXT;

-- Create unique index on firebase_id for fast lookups
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_firebase_id ON users(firebase_id);
