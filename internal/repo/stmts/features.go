package stmts

var (
	CreateFeatureTable = Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS features (
			name TEXT PRIMARY KEY UNIQUE NOT NULL,
			description TEXT,
		)`,
	}

	CreateUserFeatureTable = Statement[any]{
		Query: `CREATE TABLE IF NOT EXISTS user_features (
			user_id TEXT NOT NULL,
			feature_name TEXT NOT NULL,
			CONSTRAINT fk_feature_user_id
				FOREIGN KEY (user_id) REFERENCES(users.id)
				ON DELETE CASCADE,
			CONSTRAINT fk_feature_feature_name
				FOREIGN KEY (feature_name) REFERENCES(features.name)
				ON DELETE CASCADE
		)`,
	}

	CreateFeature = Statement[any]{
		Query: `INSERT INTO features SET name = ?, description = ?`,
		Args:  []string{"name", "description"},
	}

	EnableFeature = Statement[any]{
		Query: `INSERT INTO user_features SET feature_name = ?, user_id = ?`,
		Args:  []string{"feature_name", "user_id"},
	}

	DisableFeature = Statement[any]{
		Query: `DELETE FROM user_features WHERE feature_name = ? AND user_id = ?`,
		Args:  []string{"feature_name", "user_id"},
	}

	GetEnabledFeatures = Statement[any]{
		Query: `SELECT * from features
			JOIN user_features ON feature_name = name
			WHERE name = ? AND user_id = ?`,
		Args: []string{"name", "user_id"},
	}
)
