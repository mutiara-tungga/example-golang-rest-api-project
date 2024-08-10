BEGIN;
  ALTER TABLE users
    ADD CONSTRAINT user_unique_username UNIQUE (username);
END;