-- +migrate Up
CREATE TABLE IF NOT EXISTS roles (
    id TEXT NOT NULL PRIMARY KEY UNIQUE,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    delete_protected BOOL NOT NULL DEFAULT FALSE,
    origin TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS role_assignments (
    user_id TEXT NOT NULL,
    role_id TEXT NOT NULL,
    CONSTRAINT fk_user_role_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_role_role FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE,
    UNIQUE(user_id, role_id)
);

-- +migrate Down
DROP TABLE roles;
DROP TABLE role_assignments;