ALTER TABLE media
  RENAME COLUMN unprocessed_file_url TO file_url;

ALTER TABLE media
  RENAME CONSTRAINT media_unprocessed_file_url_key TO media_file_url_key;
