CREATE TABLE accounts (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  password_algorithm TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE temporary_accounts (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  number_code INTEGER NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE auth_audit_logs (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  success BOOLEAN NOT NULL,
  account_id UUID, -- 存在しないアカウントに試行する場合は NULL
  event_type TEXT NOT NULL, -- SIGN_UP / LOGIN / LOGOUT
  identifier TEXT, -- email 等（失敗時用）
  failure_reason TEXT,
  ip_address INET NOT NULL,
  user_agent TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);