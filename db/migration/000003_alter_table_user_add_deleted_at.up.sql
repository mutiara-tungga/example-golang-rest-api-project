BEGIN;
  ALTER table users 
    ADD COLUMN deleted_by varchar(255) NULL,
    ADD COLUMN deleted_at timestamptz NULL; 
END;