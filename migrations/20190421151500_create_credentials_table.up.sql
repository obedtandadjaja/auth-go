CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SEQUENCE credentials_id_seq;
CREATE TABLE credentials (
  id             int PRIMARY KEY nextval('credentials_id_seq'),
  uuid           uuid DEFAULT uuid_generate_v4(),
  password       varchar(128),
  last_signed_in timestamp,
  created_at     timestamp,
  updated_at     timestamp
);
ALTER SEQUENCE credentials_id_seq OWNED BY credentials.id;
CREATE INDEX credentials_uuid_idx ON credentials(uuid);
