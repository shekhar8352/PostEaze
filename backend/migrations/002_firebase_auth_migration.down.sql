-- Rollback migration for Firebase authentication changes
-- This will restore the schema to its original state before 003

-- Step 1: Drop platforms column
ALTER TABLE users DROP COLUMN IF EXISTS platforms;

-- Step 2: Make email column NOT NULL again
ALTER TABLE users ALTER COLUMN email SET NOT NULL;

-- Step 3: Re-add password column
ALTER TABLE users ADD COLUMN IF NOT EXISTS password TEXT;

-- Step 4: Re-add user_type column
ALTER TABLE users ADD COLUMN IF NOT EXISTS user_type TEXT;

-- Step 5: Drop platforms index
DROP INDEX IF EXISTS idx_users_platforms;
