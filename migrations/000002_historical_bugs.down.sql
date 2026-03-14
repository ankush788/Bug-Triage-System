-- Drop the historical_bugs table and remove extension if needed
DROP EXTENSION IF EXISTS vector;
DROP TABLE IF EXISTS historical_bugs;

-- Note: dropping the extension may fail if other tables use it.  It
-- can be left in place and is harmless once installed.
-- DROP EXTENSION IF EXISTS vector;
