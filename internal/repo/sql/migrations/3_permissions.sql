-- +migrate Up
CREATE TABLE IF NOT EXISTS permissions (
    name TEXT NOT NULL PRIMARY KEY,
    description TEXT
);

CREATE TABLE IF NOT EXISTS role_permissions (
    permission TEXT NOT NULL,
    role_id TEXT NOT NULL,
    CONSTRAINT fk_role_permissions_permission FOREIGN KEY(permission) REFERENCES permissions(name) ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_role FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE permission;
DROP TABLE role_permissions;