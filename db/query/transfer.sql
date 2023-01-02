

-- name: CreateTransfer :one
INSERT INTO transfer (
from_account_id,
to_account_id,
amount) VALUES ($1, $2, $3)
RETURNING *;


-- name: GetTransfer :one
SELECT * FROM transfer
WHERE id = $1 limit 1;

-- name: ListTransfers :many
SELECT * FROM transfer
ORDER BY id
limit $1
offset $2;

-- name: ListTransfersFromAccountId :many
SELECT * FROM transfer
where from_account_id = $1
ORDER BY id
limit $2
offset $3;


-- name: DeleteTransfer :exec
DELETE FROM transfer
WHERE id = $1;