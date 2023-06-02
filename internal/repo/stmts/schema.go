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
			password TEXT NOT NULL
		)`,
	}

	CreateAddressTable = Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS user_addresses (
			user_id TEXT NOT NULL,
			city_code TEXT,
			city_name TEXT,
			street TEXT,
			extra TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
	}

	CreateEMailTable = Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS user_emails (
			user_id TEXT NOT NULL,
			address TEXT NOT NULL,
			verified BOOL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
	}

	CreatePhoneNumberTable = Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS user_phone_numbers (
			user_id TEXT NOT NULL,
			phone_number TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
	}

	CreateGroupTable = Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS groups (
			id TEXT NOT NULL PRIMARY KEY UNIQUE,
			name TEXT NOT NULL UNIQUE,
			description TEXT
		)`,
	}

	CreateGroupMembershipTable = Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS group_memberships (
			user_id TEXT NOT NULL,
			group_id TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(group_id) REFERENCES groups(id)
		)`,
	}
)
