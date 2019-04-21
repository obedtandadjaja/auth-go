CREATE TABLE credentials (
  id             integer PRIMARY KEY,
  identifier     varchar(255) NOT NULL,
  password       varchar(128),
  subject        varchar(100) NOT NULL,
  last_signed_in timestamp,
  created_at     timestamp,
  updated_at     timestamp
);