-- name: EnrollUserTOTPSecret :exec
UPDATE users SET totp_secret = ? WHERE id = ? AND totp_secret IS NULL;

-- name: RemoveUserTOTPSecret :execrows
UPDATE users SET totp_secret = NULL WHERE id = ? AND totp_secret IS NOT NULL;

-- name: SetUserPassword :execrows
UPDATE users SET password = ? WHERE id = ?;

-- name: UserHasTOTPEnrolled :one
SELECT totp_secret IS NOT NULL FROM users WHERE id = ?;