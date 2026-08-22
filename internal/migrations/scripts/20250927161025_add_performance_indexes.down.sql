-- Remove performance indexes

-- Credentials table indexes
DROP INDEX IF EXISTS idx_credentials_username;
DROP INDEX IF EXISTS idx_credentials_user_id;
DROP INDEX IF EXISTS idx_credentials_username_active;

-- Sessions table indexes
DROP INDEX IF EXISTS idx_sessions_token;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_sessions_created_at;

-- Password resets table indexes
DROP INDEX IF EXISTS idx_password_resets_ott;
DROP INDEX IF EXISTS idx_password_resets_credential_id;

-- Refresh tokens table indexes
DROP INDEX IF EXISTS idx_refresh_tokens_token;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP INDEX IF EXISTS idx_refresh_tokens_expiry;