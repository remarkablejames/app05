-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE session_status AS ENUM ('active', 'expired', 'revoked');

CREATE TABLE sessions (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          token VARCHAR(255) NOT NULL UNIQUE,
                          refresh_token VARCHAR(255) NOT NULL UNIQUE,
                          status session_status NOT NULL DEFAULT 'active',
                          device_info JSONB NOT NULL DEFAULT '{}',
                          expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
                          last_activity_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
                          revoked_at TIMESTAMP WITH TIME ZONE,
                          revoked_reason TEXT,

                          CONSTRAINT valid_expiry CHECK (expires_at > created_at),
                          CONSTRAINT valid_revocation CHECK (
                              (status != 'revoked') OR
                              (revoked_at IS NOT NULL AND revoked_reason IS NOT NULL)
                              )
);

-- Indexes for better query performance
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_status ON sessions(status);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Index for active session lookup
CREATE INDEX idx_active_sessions ON sessions(user_id, status)
    WHERE status = 'active';

-- Function to automatically update last_activity_at
CREATE OR REPLACE FUNCTION update_last_activity_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_activity_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_session_activity
    BEFORE UPDATE ON sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_last_activity_at();

-- Function to automatically expire sessions
CREATE OR REPLACE FUNCTION expire_old_sessions()
RETURNS TRIGGER AS $$
BEGIN
    -- Only mark session expired if it has expired and isn't already expired.
    IF NEW.expires_at < CURRENT_TIMESTAMP AND NEW.status <> 'expired' THEN
        NEW.status := 'expired';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER auto_expire_sessions
    AFTER INSERT OR UPDATE ON sessions
                        EXECUTE FUNCTION expire_old_sessions();