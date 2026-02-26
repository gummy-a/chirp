CREATE TABLE media (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  owner_account_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  file_type TEXT NOT NULL, -- image / video
  original_file_name TEXT NOT NULL, -- uriにする前の元ファイル名
  unprocessed_file_url TEXT NOT NULL UNIQUE, -- 未加工の元ファイルを格納したurl
  metadata JSONB NOT NULL DEFAULT '{}'::JSONB -- thumbnail, 圧縮済みファイルurl, サイズ等
);
