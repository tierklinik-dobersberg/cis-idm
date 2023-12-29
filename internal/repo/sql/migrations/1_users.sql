-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id TEXT NOT NULL PRIMARY KEY UNIQUE,
    username TEXT NOT NULL UNIQUE,
    display_name TEXT NOT NULL DEFAULT '',
    first_name TEXT NOT NULL DEFAULT '',
    last_name TEXT NOT NULL DEFAULT '',
    extra TEXT NOT NULL DEFAULT '',
    avatar TEXT NOT NULL DEFAULT '',
    birthday TEXT NOT NULL DEFAULT '',
    password TEXT NOT NULL,
    totp_secret TEXT 
);

CREATE TABLE IF NOT EXISTS user_addresses (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,
    city_code TEXT NOT NULL DEFAULT '',
    city_name TEXT NOT NULL DEFAULT '',
    street TEXT NOT NULL DEFAULT '',
    extra TEXT NOT NULL DEFAULT '',
    CONSTRAINT fk_user_address FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_emails (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,
    address TEXT UNIQUE NOT NULL,
    verified BOOL NOT NULL DEFAULT FALSE,
    is_primary BOOL NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_user_mail FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS user_phone_numbers (
    id TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    is_primary BOOL NOT NULL DEFAULT FALSE,
    verified BOOL NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_user_phone_number FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS mfa_backup_codes (
    code TEXT NOT NULL,
    user_id TEXT NOT NULL,
    CONSTRAINT fk_mfa_backup_codes_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS webauthn_creds (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    cred TEXT NOT NULL,
    cred_type TEXT NOT NULL,
    client_name TEXT NOT NULL,
    client_os TEXT NOT NULL,
    client_device TEXT NOT NULL,
    CONSTRAINT fk_webauth_creds_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS webpush_subscriptions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    user_agent TEXT NOT NULL DEFAULT '',
    endpoint TEXT NOT NULL UNIQUE,
    auth TEXT NOT NULL,
    key TEXT NOT NULL,
    token_id TEXT NOT NULL,
    CONSTRAINT fk_webpush_subscription_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token_invalidation (
    token_id TEXT NOT NULL PRIMARY KEY UNIQUE,
    user_id TEXT NOT NULL,
    issued_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_token_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE users;
DROP TABLE user_addresses;
DROP TABLE user_emails;
DROP TABLE user_phone_numbers;
DROP TABLE mfa_backup_codes;
DROP TABLE webauthn_creds;
DROP TABLE webpush_subscriptions;
DROP TABLE token_invalidation;