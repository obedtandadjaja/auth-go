CREATE SEQUENCE IF NOT EXISTS refresh_tokens_id_seq;
CREATE TABLE IF NOT EXISTS refresh_tokens(
  id               int PRIMARY KEY DEFAULT nextval('refresh_tokens_id_seq'),
  uuid             uuid DEFAULT uuid_generate_val(),
  credential_id    int REFERENCES credentials(id),
  ip_address       varchar(100),
  user_agent       varchar(255),
  last_accessed_at timestamp DEFAULT now(),
  created_at       timestamp DEFAULT now(),
  expires_at       timestamp DEFAULT now()
);
ALTER SEQUENCE refresh_tokens_id_seq OWNED BY refresh_tokens.id;
CREATE INDEX IF NOT EXISTS refresh_tokens_token_idx ON refresh_tokens(token);
