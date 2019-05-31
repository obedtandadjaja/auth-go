ALTER TABLE credentials ADD COLUMN failed_attempts integer DEFAULT 0;
ALTER TABLE credentials ADD COLUMN locked_until timestamp;
