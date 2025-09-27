-- Add performance indexes for common query patterns

-- Credentials table indexes
CREATE INDEX idx_credentials_username ON credentials(username);
CREATE INDEX idx_credentials_user_id ON credentials(user_id);
CREATE INDEX idx_credentials_username_active ON credentials(username, active);

-- Sessions table indexes
CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_created_at ON sessions(created_at);

-- Password resets table indexes
CREATE INDEX idx_password_resets_ott ON password_resets(ott);
CREATE INDEX idx_password_resets_credential_id ON password_resets(credential_id);

-- Refresh tokens table indexes
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expiry ON refresh_tokens(expiry);