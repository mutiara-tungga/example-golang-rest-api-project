BEGIN;
  ALTER TABLE users
    DROP CONSTRAINT IF EXISTS user_unique_username;
END;