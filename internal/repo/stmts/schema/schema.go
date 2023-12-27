package schema

import "github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"

var createSchema = stmts.StatementList{
	// Schema table
	stmts.Statement[any]{
		Query: `CREATE TABLE schema_version (
			version TEXT PRIMARY KEY
		)`,
	},

	// Users
	stmts.Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS users (
			id TEXT NOT NULL PRIMARY KEY UNIQUE,
			username TEXT NOT NULL UNIQUE,
			display_name TEXT,
			first_name TEXT,
			last_name TEXT,
			extra TEXT,
			avatar TEXT,
			birthday TEXT,
			password TEXT NOT NULL,
			totp_secret TEXT NULL
		)`,
	},

	// Roles
	stmts.Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS roles (
			id TEXT NOT NULL PRIMARY KEY UNIQUE,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			delete_protected BOOL
		)`,
	},

	// Permissions
	stmts.Statement[any] {
		Query: `CREATE TABLE IF NOT EXISTS permissions (
			name TEXT NOT NULL PRIMARY KEY,
			description TEXT
		)`,
	},

	// API tokens
	stmts.Statement[any] {
		Query: `CREATE TABLE IF NOT EXISTS user_api_tokens (
			token TEXT NOT NULL PRIMARY KEY UNIQUE,
			name TEXT NOT NULL,
			user_id TEXT NOT NULL,
			expires_at NUMBER,
			created_at NUMBER NOT NULL,
			CONSTRAINT fk_user_api_token 
				FOREIGN KEY(user_id) REFERENCES users(id)
				ON DELETE CASCADE
		)`,
	},

	stmts.Statement[any] {
		Query: `CREATE TABLE IF NOT EXISTS user_api_token_roles (
			token TEXT NOT NULL,
			role_id TEXT NOT NULL,
			CONSTRAINT fk_user_api_token_roles_token
				FOREIGN KEY(token) REFERENCES user_api_tokens(token)
				ON DELETE CASCADE,
			CONSTRAINT fk_user_api_token_roles_role
				FOREIGN KEY(role_id) REFERENCES roles(id)
				ON DELETE CASCADE
		)`,
	},

	stmts.Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS role_permissions (
			permission TEXT NOT NULL,
			role_id TEXT NOT NULL,
			CONSTRAINT fk_role_permissions_permission 
				FOREIGN KEY(permission) REFERENCES permissions(name)
				ON DELETE CASCADE,
			CONSTRAINT fk_role_permissions_role
				FOREIGN KEY(role_id) REFERNCES roles(id)
				ON DELETE CASCADE
		)`,
	},

	stmts.Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS user_addresses (
			id TEXT NOT NULL PRIMARY KEY,
			user_id TEXT NOT NULL,
			city_code TEXT,
			city_name TEXT,
			street TEXT,
			extra TEXT,
			CONSTRAINT fk_user_address
				FOREIGN KEY(user_id) REFERENCES users(id)
				ON DELETE CASCADE
		)`,
	},

	stmts.Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS user_emails (
			id TEXT NOT NULL PRIMARY KEY,
			user_id TEXT NOT NULL,
			address TEXT UNIQUE NOT NULL,
			verified BOOL,
			is_primary BOOL,
			CONSTRAINT fk_user_mail
				FOREIGN KEY(user_id) REFERENCES users(id)
				ON DELETE CASCADE
		)`,
	},

	stmts.Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS user_phone_numbers (
			id TEXT NOT NULL PRIMARY KEY,
			user_id TEXT NOT NULL,
			phone_number TEXT NOT NULL,
			is_primary BOOL,
			verified BOOL,
			CONSTRAINT fk_user_phone_number
				FOREIGN KEY(user_id) REFERENCES users(id)
				ON DELETE CASCADE
		)`,
	},

	stmts.Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS role_assignments (
			user_id TEXT NOT NULL,
			role_id TEXT NOT NULL,
			CONSTRAINT fk_user_role_user
				FOREIGN KEY(user_id) REFERENCES users(id)
				ON DELETE CASCADE,
			CONSTRAINT fk_user_role_role
				FOREIGN KEY(role_id) REFERENCES roles(id)
				ON DELETE CASCADE
		)`,
	},

	stmts.Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS token_invalidation (
			token_id TEXT NOT NULL PRIMARY KEY UNIQUE,
			user_id TEXT NOT NULL,
			issued_at NUMBER NOT NULL,
			expires_at NUMBER NOT NULL,
			CONSTRAINT fk_token_user
				FOREIGN KEY(user_id) REFERENCES users(id)
				ON DELETE CASCADE
		)`,
	},

	stmts.Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS registration_tokens (
			token TEXT NOT NULL PRIMARY KEY,
			expires NUMBER NULL,
			allowed_usage NUMBER NULL,
			initial_roles TEXT NULL,
			created_by STRING,
			created_at NUMBER
		)`,
	},

	stmts.Statement[any]{
		Query: `CREATE TRIGGER IF NOT EXISTS registration_token_cleanup AFTER UPDATE ON registration_tokens
		BEGIN
			DELETE FROM registration_tokens
			WHERE
				allowed_usage IS NOT NULL
				AND allowed_usage = 0;
		END;`,
	},

	stmts.Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS mfa_backup_codes (
			code TEXT NOT NULL,
			user_id TEXT,
			CONSTRAINT fk_mfa_backup_codes_user
				FOREIGN KEY(user_id) REFERENCES users(id)
				ON DELETE CASCADE
		)`,
	},

	stmts.Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS webauthn_creds (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			cred TEXT,
			cred_type TEXT,
			client_name TEXT,
			client_os TEXT,
			client_device TEXT,
			CONSTRAINT fk_webauth_creds_user
				FOREIGN KEY(user_id) REFERENCES users(id)
				ON DELETE CASCADE
		)`,
	},

	stmts.Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS webpush_subscriptions (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			user_agent TEXT,
			endpoint TEXT NOT NULL UNIQUE,
			auth TEXT NOT NULL,
			key TEXT NOT NULL,
			token_id TEXT NOT NULL,
			CONSTRAINT fk_webpush_subscription_user
				FOREIGN KEY(user_id) REFERENCES users(id)
				ON DELETE CASCADE
		)
		`,
	},
}
