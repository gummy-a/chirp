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
