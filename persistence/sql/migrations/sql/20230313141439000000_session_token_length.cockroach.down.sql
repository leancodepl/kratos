-- Downsizing is not yet supported in CockroachDB. Since this migration has no real-world impact on the application, we can safely
-- not execute it.
--
-- ALTER TABLE sessions ALTER COLUMN token TYPE varchar(32);
-- ALTER TABLE sessions ALTER COLUMN logout_token TYPE varchar(32);