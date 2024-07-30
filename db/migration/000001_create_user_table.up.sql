BEGIN;
  CREATE TABLE users(
      id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
      name varchar(255) NOT NULL,
      username varchar(100) NOT NULL,
      phone varchar(15) NOT NULL DEFAULT '',
      password text NOT NULL,
      created_at timestamptz NOT NULL DEFAULT NOW(),
      created_by varchar(255),
      updated_at timestamptz NOT NULL DEFAULT NOW(),
      updated_by varchar(255)
  );
END;