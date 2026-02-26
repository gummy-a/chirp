-- name: SelectMediaById :one
SELECT *
FROM media
WHERE id = $1;

-- name: SelectMediaByOwnerAccountId :many
SELECT *
FROM media
WHERE owner_account_id = $1
LIMIT $2 OFFSET $3;

-- name: InsertMedia :one
INSERT INTO media (owner_account_id, file_type, unprocessed_file_url, metadata)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateMedia :one
UPDATE media
SET metadata = $2
WHERE id = $1
RETURNING *;

-- name: DeleteMedia :exec
DELETE FROM media
WHERE id = $1;

-- name: DeleteMediaByOwnerId :exec
DELETE FROM media
WHERE owner_account_id = $1;
