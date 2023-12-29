-- +migrate Up
CREATE TABLE IF NOT EXISTS user_api_tokens (
    token TEXT NOT NULL PRIMARY KEY UNIQUE,
    name TEXT NOT NULL,
    user_id TEXT NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_api_token FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_api_token_roles (
    token TEXT NOT NULL,
    role_id TEXT NOT NULL,
    CONSTRAINT fk_user_api_token_roles_token FOREIGN KEY(token) REFERENCES user_api_tokens(token) ON DELETE CASCADE,
    CONSTRAINT fk_user_api_token_roles_role FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE user_api_tokens;
DROP TABLE user_api_token_roles;