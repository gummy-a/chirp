ALTER TABLE media
  RENAME CONSTRAINT media_file_url_key TO media_unprocessed_file_url_key;

ALTER TABLE media
  RENAME COLUMN file_url TO unprocessed_file_url;
