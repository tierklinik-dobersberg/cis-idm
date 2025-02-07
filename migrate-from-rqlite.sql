.bail on

PRAGMA foreign_keys = ON;

ATTACH DATABASE 'file:///tmp/backup.db' AS backup;

BEGIN TRANSACTION;

    -- Update users table
    UPDATE backup.users SET display_name = '' WHERE display_name IS NULL;
    UPDATE backup.users SET first_name = '' WHERE first_name IS NULL;
    UPDATE backup.users SET last_name = '' WHERE last_name IS NULL;
    UPDATE backup.users SET extra = '' WHERE extra IS NULL;
    UPDATE backup.users SET avatar = '' WHERE avatar IS NULL;
    UPDATE backup.users SET birthday = '' WHERE birthday IS NULL;

    -- Update user_addresses table
    UPDATE backup.user_addresses SET city_code = '' WHERE city_code IS NULL;
    UPDATE backup.user_addresses SET city_name = '' WHERE city_name IS NULL;
    UPDATE backup.user_addresses SET street = '' WHERE street IS NULL;
    UPDATE backup.user_addresses SET extra = '' WHERE extra IS NULL;

    -- Update user_emails table
    UPDATE backup.user_emails SET verified = FALSE WHERE verified IS NULL;
    UPDATE backup.user_emails SET is_primary = FALSE WHERE is_primary IS NULL;

    -- Update user_phone_numbers table
    UPDATE backup.user_phone_numbers SET verified = FALSE WHERE verified IS NULL;
    UPDATE backup.user_phone_numbers SET is_primary = FALSE WHERE is_primary IS NULL;

    -- Update webpush subscriptions
    UPDATE backup.webpush_subscriptions SET user_agent = '' WHERE user_agent IS NULL;

    -- Update roles
    UPDATE backup.roles SET description = '' WHERE description IS NULL;
    UPDATE backup.roles SET delete_protected = FALSE where delete_protected IS NULL;

    -- Start copying
    INSERT OR ROLLBACK INTO users SELECT * from backup.users;
    INSERT OR ROLLBACK INTO user_addresses SELECT * from backup.user_addresses;
    INSERT OR ROLLBACK INTO user_emails SELECT * from backup.user_emails;
    INSERT OR ROLLBACK INTO user_phone_numbers SELECT * from backup.user_phone_numbers;
    INSERT OR ROLLBACK INTO mfa_backup_codes SELECT * from backup.mfa_backup_codes;
    INSERT OR ROLLBACK INTO webauthn_creds SELECT * from backup.webauthn_creds;
    INSERT OR ROLLBACK INTO webpush_subscriptions SELECT * from backup.webpush_subscriptions;
    INSERT OR ROLLBACK INTO token_invalidation SELECT * from backup.token_invalidation;
    
    INSERT OR ROLLBACK INTO roles SELECT *, 'api' from backup.roles;
    INSERT OR ROLLBACK INTO role_assignments SELECT DISTINCT * from backup.role_assignments;

COMMIT;
