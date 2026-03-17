-- Migration to update users table for Firebase authentication
-- This migration removes password and user_type fields, makes email optional, and adds platforms array

-- Step 1: Add platforms column
ALTER TABLE users ADD COLUMN platforms TEXT[] DEFAULT ARRAY[]::TEXT[];

-- Step 2: Make email column nullable
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

-- Step 3: Drop password column (if exists)
ALTER TABLE users DROP COLUMN IF EXISTS password;

-- Step 4: Drop user_type column (if exists)
ALTER TABLE users DROP COLUMN IF EXISTS user_type;

-- Step 5: Create index for efficient platform queries
CREATE INDEX IF NOT EXISTS idx_users_platforms ON users USING GIN(platforms);

-- Step 6: Update existing users to have empty platforms array if NULL
UPDATE users SET platforms = ARRAY[]::TEXT[] WHERE platforms IS NULL;