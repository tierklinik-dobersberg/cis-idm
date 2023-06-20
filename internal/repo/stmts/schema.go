package stmts

var (
	CreateUserTable = Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS users (
			id TEXT NOT NULL PRIMARY KEY UNIQUE,
			username TEXT NOT NULL UNIQUE,
			display_name TEXT,
			first_name TEXT,
			last_name TEXT,
			extra BLOB,
			avatar TEXT,
			birthday TEXT,
			password TEXT NOT NULL
		)`,
	}

	CreateAddressTable = Statement[any]{
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
	}

	CreateEMailTable = Statement[any]{
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
	}

	CreatePhoneNumberTable = Statement[any]{
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
	}

	CreateRoleTable = Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS roles (
			id TEXT NOT NULL PRIMARY KEY UNIQUE,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			delete_protected BOOL
		)`,
	}

	CreateRoleAssignmentTable = Statement[any]{
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
	}

	CreateTokenInvalidationTable = Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS token_invalidation (
			token_id TEXT NOT NULL PRIMARY KEY UNIQUE,
			user_id TEXT NOT NULL,
			issued_at NUMBER NOT NULL,
			expires_at NUMBER NOT NULL,
			CONSTRAINT fk_token_user
				FOREIGN KEY(user_id) REFERENCES users(id)
				ON DELETE CASCADE
		)`,
	}
)
