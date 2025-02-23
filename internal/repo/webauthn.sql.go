// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: webauthn.sql

package repo

import (
	"context"
)

const addWebauthnCred = `-- name: AddWebauthnCred :exec
INSERT INTO webauthn_creds (id, user_id, cred, client_name, client_os, client_device, cred_type) VALUES (?, ?, ?, ?, ?, ?, ?)
`

type AddWebauthnCredParams struct {
	ID           string
	UserID       string
	Cred         string
	ClientName   string
	ClientOs     string
	ClientDevice string
	CredType     string
}

func (q *Queries) AddWebauthnCred(ctx context.Context, arg AddWebauthnCredParams) error {
	_, err := q.db.ExecContext(ctx, addWebauthnCred,
		arg.ID,
		arg.UserID,
		arg.Cred,
		arg.ClientName,
		arg.ClientOs,
		arg.ClientDevice,
		arg.CredType,
	)
	return err
}

const getWebauthnCreds = `-- name: GetWebauthnCreds :many
SELECT id, user_id, cred, cred_type, client_name, client_os, client_device FROM webauthn_creds WHERE user_id = ?
`

func (q *Queries) GetWebauthnCreds(ctx context.Context, userID string) ([]WebauthnCred, error) {
	rows, err := q.db.QueryContext(ctx, getWebauthnCreds, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []WebauthnCred
	for rows.Next() {
		var i WebauthnCred
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Cred,
			&i.CredType,
			&i.ClientName,
			&i.ClientOs,
			&i.ClientDevice,
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

const removeWebauthnCred = `-- name: RemoveWebauthnCred :execrows
DELETE FROM webauthn_creds WHERE user_id = ? AND id = ?
`

type RemoveWebauthnCredParams struct {
	UserID string
	ID     string
}

func (q *Queries) RemoveWebauthnCred(ctx context.Context, arg RemoveWebauthnCredParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, removeWebauthnCred, arg.UserID, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
