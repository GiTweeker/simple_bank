
-- name: CreateEntryForAccount :one
INSERT INTO entries (
account_id,
amount) VALUES ($1, $2)
RETURNING *;


-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 limit 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id
limit $1
offset $2;

-- name: ListEntriesForAccount :many
SELECT * FROM entries
where account_id = $1
ORDER BY id
limit $2
offset $3;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;