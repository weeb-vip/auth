-- Postgres baseline for auth.
--
-- Generated from the production schema dumped out of PlanetScale, not from
-- this repository's MySQL migrations. The migrations describe what was
-- intended; the dump describes what is actually there, and on anime-api those
-- two disagreed -- a table its migrations create had been dropped from
-- production years ago and nothing noticed.
--
-- Worth knowing while reading this: every service shares one MySQL database
-- today. These tables sit alongside anime-api's, with a separate
-- __migrations_* tracker per service. Whether that stays true in Postgres is
-- an open decision; this file only covers the tables this service owns.

CREATE TABLE credentials (
  id varchar(100) NOT NULL,
  username varchar(255) NOT NULL,
  type varchar(255) NOT NULL,
  value varchar(255) NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  user_id varchar(255),
  active boolean DEFAULT false,
  PRIMARY KEY (id)
);

CREATE INDEX idx_credentials_username ON credentials (username);
CREATE INDEX idx_credentials_user_id ON credentials (user_id);
CREATE INDEX idx_credentials_username_active ON credentials (username, active);

CREATE TABLE password_resets (
  id varchar(100) NOT NULL,
  credential_id varchar(100) NOT NULL,
  ott varchar(255) NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  PRIMARY KEY (id)
);

CREATE INDEX idx_password_resets_ott ON password_resets (ott);
CREATE INDEX idx_password_resets_credential_id ON password_resets (credential_id);

CREATE TABLE refresh_tokens (
  id varchar(100) NOT NULL,
  user_id varchar(100) NOT NULL,
  token varchar(100) NOT NULL,
  expiry integer NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  PRIMARY KEY (id)
);

CREATE INDEX idx_refresh_tokens_token ON refresh_tokens (token);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
CREATE INDEX idx_refresh_tokens_expiry ON refresh_tokens (expiry);

CREATE TABLE sessions (
  id varchar(100) NOT NULL,
  user_id varchar(100) NOT NULL,
  ip_address varchar(255) NOT NULL,
  user_agent varchar(255) NOT NULL,
  token varchar(255) NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  PRIMARY KEY (id)
);

CREATE INDEX idx_sessions_token ON sessions (token);
CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_created_at ON sessions (created_at);


