-- name: AddWebauthnCred :exec
INSERT INTO webauthn_creds (id, user_id, cred, client_name, client_os, client_device, cred_type) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetWebauthnCreds :many
SELECT * FROM webauthn_creds WHERE user_id = ?;

-- name: RemoveWebauthnCred :execrows
DELETE FROM webauthn_creds WHERE user_id = ? AND id = ?;