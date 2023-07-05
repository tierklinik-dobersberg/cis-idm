package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/rqlite/gorqlite"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

type Repo struct {
	Conn *gorqlite.Connection
}

func New(endpoint string) (*Repo, error) {
	conn, err := gorqlite.Open(endpoint)
	if err != nil {
		return nil, err
	}

	return &Repo{Conn: conn}, nil
}

func (repo *Repo) Migrate(ctx context.Context) error {
	if err := stmts.CreateUserTable.Write(ctx, repo.Conn, nil); err != nil && !errors.Is(err, stmts.ErrNoRowsAffected) {
		return fmt.Errorf("failed to create user table: %w", err)
	}

	if err := stmts.CreateAddressTable.Write(ctx, repo.Conn, nil); err != nil && !errors.Is(err, stmts.ErrNoRowsAffected) {
		return fmt.Errorf("failed to create user_addresses table: %w", err)
	}

	if err := stmts.CreatePhoneNumberTable.Write(ctx, repo.Conn, nil); err != nil && !errors.Is(err, stmts.ErrNoRowsAffected) {
		return fmt.Errorf("failed to create user_phone_numbers table: %w", err)
	}

	if err := stmts.CreateEMailTable.Write(ctx, repo.Conn, nil); err != nil && !errors.Is(err, stmts.ErrNoRowsAffected) {
		return fmt.Errorf("failed to create user_emails table: %w", err)
	}

	if err := stmts.CreateRoleTable.Write(ctx, repo.Conn, nil); err != nil && !errors.Is(err, stmts.ErrNoRowsAffected) {
		return fmt.Errorf("failed to create roles table: %w", err)
	}

	if err := stmts.CreateRoleAssignmentTable.Write(ctx, repo.Conn, nil); err != nil && !errors.Is(err, stmts.ErrNoRowsAffected) {
		return fmt.Errorf("failed to create role_assignments table: %w", err)
	}

	if err := stmts.CreateTokenInvalidationTable.Write(ctx, repo.Conn, nil); err != nil && !errors.Is(err, stmts.ErrNoRowsAffected) {
		return fmt.Errorf("failed to create token invalidation table: %w", err)
	}

	if err := stmts.CreateRegistrationTokenTable.Write(ctx, repo.Conn, nil); err != nil && !errors.Is(err, stmts.ErrNoRowsAffected) {
		return fmt.Errorf("failed to create registration_token table: %w", err)
	}

	if err := stmts.CreateRegistrationTokenCleanupTrigger.Write(ctx, repo.Conn, nil); err != nil && !errors.Is(err, stmts.ErrNoRowsAffected) {
		return fmt.Errorf("failed to create registration_token_cleanup trigger: %w", err)
	}

	if err := stmts.Create2FABackupCodeTable.Write(ctx, repo.Conn, nil); err != nil && !errors.Is(err, stmts.ErrNoRowsAffected) {
		return fmt.Errorf("failed to create mfa_backup_codes table: %w", err)
	}

	if err := stmts.CreateWebauthnCredsTable.Write(ctx, repo.Conn, nil); err != nil && !errors.Is(err, stmts.ErrNoRowsAffected) {
		return fmt.Errorf("failed to create webauthn_creds table: %w", err)
	}

	if err := stmts.CreateWebauthnSessionTable.Write(ctx, repo.Conn, nil); err != nil && !errors.Is(err, stmts.ErrNoRowsAffected) {
		return fmt.Errorf("failed to create webauthn_sessions table: %w", err)
	}

	return nil
}
