BEGIN;
  ALTER table users 
    DROP COLUMN deleted_by,
    DROP COLUMN deleted_at; 
END;