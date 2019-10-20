CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE credentials (
  id             uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  password       varchar(128),
  last_signed_in timestamp,
  created_at     timestamp,
  updated_at     timestamp
);
