-- name: GetWebPushSubscriptionsForUser :many
SELECT
	*
FROM
	webpush_subscriptions
WHERE
	user_id = ?;

-- name: DeleteWebPushSubscriptionForToken :execrows
DELETE FROM
	webpush_subscriptions
WHERE
	token_id = ?;

-- name: DeleteWebPushSubscriptionByID :execrows
DELETE FROM
	webpush_subscriptions
WHERE
	id = ?;

-- name: CreateWebPushSubscriptionForUser :exec
INSERT
	OR REPLACE INTO webpush_subscriptions (
		id,
		user_id,
		user_agent,
		endpoint,
		auth,
		key,
		token_id
	)
VALUES
	(?, ?, ?, ?, ?, ?, ?);