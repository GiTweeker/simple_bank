// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: entries.sql

package db

import (
	"context"
)

const createEntryForAccount = `-- name: CreateEntryForAccount :one
INSERT INTO entries (
account_id,
amount) VALUES ($1, $2)
RETURNING id, account_id, amount, created_at
`

type CreateEntryForAccountParams struct {
	AccountID int64 `json:"accountID"`
	Amount    int64 `json:"amount"`
}

func (q *Queries) CreateEntryForAccount(ctx context.Context, arg CreateEntryForAccountParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntryForAccount, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const deleteEntry = `-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1
`

func (q *Queries) DeleteEntry(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteEntry, id)
	return err
}

const getEntry = `-- name: GetEntry :one
SELECT id, account_id, amount, created_at FROM entries
WHERE id = $1 limit 1
`

func (q *Queries) GetEntry(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntry, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listEntries = `-- name: ListEntries :many
SELECT id, account_id, amount, created_at FROM entries
ORDER BY id
limit $1
offset $2
`

type ListEntriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, listEntries, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listEntriesForAccount = `-- name: ListEntriesForAccount :many
SELECT id, account_id, amount, created_at FROM entries
where account_id = $1
ORDER BY id
limit $2
offset $3
`

type ListEntriesForAccountParams struct {
	AccountID int64 `json:"accountID"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListEntriesForAccount(ctx context.Context, arg ListEntriesForAccountParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, listEntriesForAccount, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
