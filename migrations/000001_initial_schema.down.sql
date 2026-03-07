-- Drop indexes
DROP INDEX IF EXISTS idx_bugs_priority;
DROP INDEX IF EXISTS idx_bugs_status;
DROP INDEX IF EXISTS idx_bugs_reporter_id;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables
DROP TABLE IF EXISTS bugs;
DROP TABLE IF EXISTS users;