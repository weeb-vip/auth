-- Remove performance indexes

-- Credentials table indexes
DROP INDEX idx_credentials_username ON credentials;
DROP INDEX idx_credentials_user_id ON credentials;
DROP INDEX idx_credentials_username_active ON credentials;

-- Sessions table indexes
DROP INDEX idx_sessions_token ON sessions;
DROP INDEX idx_sessions_user_id ON sessions;
DROP INDEX idx_sessions_created_at ON sessions;

-- Password resets table indexes
DROP INDEX idx_password_resets_ott ON password_resets;
DROP INDEX idx_password_resets_credential_id ON password_resets;

-- Refresh tokens table indexes
DROP INDEX idx_refresh_tokens_token ON refresh_tokens;
DROP INDEX idx_refresh_tokens_user_id ON refresh_tokens;
DROP INDEX idx_refresh_tokens_expiry ON refresh_tokens;