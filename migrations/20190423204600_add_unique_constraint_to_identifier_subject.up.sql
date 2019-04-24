ALTER TABLE credentials ADD CONSTRAINT credentials_identifier_subject_idx UNIQUE (identifier, subject);
