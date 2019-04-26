ALTER TABLE credentials
  ADD COLUMN failed_attempts integer
  ADD COLUMN locked_until timestamp;
