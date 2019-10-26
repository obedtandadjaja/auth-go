CREATE SEQUENCE refresh_tokens_id_seq;
CREATE TABLE refresh_tokens (
  id             int PRIMARY KEY nextval('refresh_tokens_id_seq'),
  token          varchar(1000) NOT NULL,
  credentials_id int REFERENCES credentials(id),
  expires_at     timestamp
);
ALTER SEQUENCE refresh_tokens_id_seq OWNED BY refresh_tokens.id;
CREATE INDEX refresh_tokens_token_idx ON refresh_tokens(token);
