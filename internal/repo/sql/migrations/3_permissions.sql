-- +migrate Up
CREATE TABLE IF NOT EXISTS role_permissions (
    permission TEXT NOT NULL,
    role_id TEXT NOT NULL,
    CONSTRAINT fk_role_permissions_role FOREIGN KEY(role_id) REFERENCES roles(id) ON DELETE CASCADE,
    UNIQUE(permission, role_id)
);

-- +migrate Down
DROP TABLE role_permissions;