DROP TRIGGER IF EXISTS auto_expire_sessions ON sessions;
DROP FUNCTION IF EXISTS expire_old_sessions();
DROP TRIGGER IF EXISTS update_session_activity ON sessions;
DROP FUNCTION IF EXISTS update_last_activity_at();
DROP TABLE IF EXISTS sessions;
DROP TYPE IF EXISTS session_status;