-- name: CreateAccount :one
INSERT INTO Accounts (
owner,
balance,
currency) VALUES ($1, $2, $3)
RETURNING *;


-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 limit 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 limit 1
FOR no key update ;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
limit $1
offset $2;

-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id  = sqlc.arg(id)
returning *;


-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id  = $1
returning *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;