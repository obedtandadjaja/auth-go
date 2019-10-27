CREATE SEQUENCE IF NOT EXISTS refresh_tokens_id_seq;
CREATE TABLE IF NOT EXISTS refresh_tokens(
  id             int PRIMARY KEY DEFAULT nextval('refresh_tokens_id_seq'),
  uuid           uuid DEFAULT uuid_generate_val(),
  credential_id  int REFERENCES credentials(id),
  expires_at     timestamp
);
ALTER SEQUENCE refresh_tokens_id_seq OWNED BY refresh_tokens.id;
CREATE INDEX IF NOT EXISTS refresh_tokens_token_idx ON refresh_tokens(token);
