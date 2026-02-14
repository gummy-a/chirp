-- name: CreateAccount :one
INSERT INTO accounts (email, password_hash, password_algorithm)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteAccount :one
DELETE FROM accounts
WHERE id = $1
RETURNING *;

-- name: FindAccountByEmail :one
SELECT *
FROM accounts
WHERE email = $1;

-- name: FindAccountById :one
SELECT *
FROM accounts
WHERE id = $1;


-- name: CreateTemporaryAccount :one
INSERT INTO temporary_accounts (email, password_hash, expires_at, number_code)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteTemporaryAccount :one
DELETE FROM temporary_accounts
WHERE id = $1
RETURNING *;

-- name: FindTemporaryAccountById :one
SELECT *
FROM temporary_accounts
WHERE id = $1;
