// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: totp.sql

package repo

import (
	"context"
	"database/sql"
)

const enrollUserTOTPSecret = `-- name: EnrollUserTOTPSecret :exec
UPDATE users SET totp_secret = ? WHERE id = ? AND totp_secret IS NULL
`

type EnrollUserTOTPSecretParams struct {
	TotpSecret sql.NullString
	ID         string
}

func (q *Queries) EnrollUserTOTPSecret(ctx context.Context, arg EnrollUserTOTPSecretParams) error {
	_, err := q.db.ExecContext(ctx, enrollUserTOTPSecret, arg.TotpSecret, arg.ID)
	return err
}

const removeUserTOTPSecret = `-- name: RemoveUserTOTPSecret :execrows
UPDATE users SET totp_secret = NULL WHERE id = ? AND totp_secret IS NOT NULL
`

func (q *Queries) RemoveUserTOTPSecret(ctx context.Context, id string) (int64, error) {
	result, err := q.db.ExecContext(ctx, removeUserTOTPSecret, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const setUserPassword = `-- name: SetUserPassword :execrows
UPDATE users SET password = ? WHERE id = ?
`

type SetUserPasswordParams struct {
	Password string
	ID       string
}

func (q *Queries) SetUserPassword(ctx context.Context, arg SetUserPasswordParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, setUserPassword, arg.Password, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const userHasTOTPEnrolled = `-- name: UserHasTOTPEnrolled :one
SELECT totp_secret IS NOT NULL FROM users WHERE id = ?
`

func (q *Queries) UserHasTOTPEnrolled(ctx context.Context, id string) (bool, error) {
	row := q.db.QueryRowContext(ctx, userHasTOTPEnrolled, id)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}
