package policy

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/jwt"
	"github.com/tierklinik-dobersberg/cis-idm/internal/permission"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

// SubjectInput defines the input for rego policies under the input.subject
// path and is populated from the user performing the operation.
type SubjectInput struct {
	// ID is the unique identifier of the user.
	ID string `mapstructure:"id" json:"id"`

	// Username is the name of the user.
	// SECURITY: If cisidm is configured to allow username changes using the username
	// in rego policies is a huge security risk!
	Username string `mapstructure:"username" json:"username"`

	// Roles is a list of roles assigned to the user. Note that the permissions
	// assigned to each role are not exposed to rego policies. Use the Permissions
	// field below which contains a set of resolved permissions from all user roles.
	Roles []repo.Role `mapstructure:"roles" json:"roles"`

	// Permissions holds the resolved set of permissions this user has based on all
	// assigned roles.
	Permissions []string `mapstructure:"permissions" json:"permissions"`

	// Fields hold the additional user fields as specified in the configuration.
	Fields map[string]any `mapstructure:"fields" json:"fields"`

	// Email holds the primary email address of the user.
	Email string `mapstructure:"email" json:"email"`

	// DisplayName holds the display name of the user.
	DisplayName string `mapstructure:"display_name" json:"display_name"`

	// TokenKind reports how the access token used to perform the request was
	// obtained. Valid values are "password", "mfa" and "webauthn".
	TokenKind jwt.LoginKind `mapstructure:"token_kind" json:"token_kind"`
}

type Store interface {
	GetUserByID(context.Context, string) (repo.User, error)
	GetRolesForUser(context.Context, string) ([]repo.Role, error)
	GetRolesForToken(context.Context, string) ([]repo.Role, error)
	GetRolePermissions(context.Context, string) ([]string, error)
	GetPrimaryEmailForUserByID(context.Context, string) (repo.UserEmail, error)
}

func NewSubjectInput(ctx context.Context, ds Store, permissionResolver permission.Resolver, userID string, tokenKind jwt.LoginKind, tokenID string) (*SubjectInput, error) {
	user, err := ds.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %q: %w", userID, err)
	}

	if user.Deleted {
		return nil, fmt.Errorf("user profile has been deleted")
	}

	var roles []repo.Role
	if tokenKind == jwt.LoginKindAPI {
		roles, err = ds.GetRolesForToken(ctx, tokenID)
	} else {
		roles, err = ds.GetRolesForUser(ctx, userID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user or token roles: %w", err)
	}

	permissionSet := make([]string, 0, len(roles))
	for _, role := range roles {
		rolePermissions, err := ds.GetRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get permissions for role %q: %w", role.ID, err)
		}

		permissionSet = append(permissionSet, rolePermissions...)
	}

	resolvedPermissions, err := permissionResolver.Resolve(permissionSet)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve user permissions: %w", err)
	}

	var mail string
	if primary, err := ds.GetPrimaryEmailForUserByID(ctx, userID); err == nil {
		mail = primary.Address
	} else {
		if !errors.Is(err, sql.ErrNoRows) {
			log.L(ctx).Error("failed to get primary email address", "user", userID, "error", err)
		}
	}

	input := &SubjectInput{
		Username:    user.Username,
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Email:       mail,
		Roles:       roles,
		Permissions: resolvedPermissions,
		TokenKind:   tokenKind,
		Fields:      make(map[string]any),
	}

	if len(user.Extra) > 0 {
		if err := json.Unmarshal([]byte(user.Extra), &input.Fields); err != nil {
			log.L(ctx).Error("failed to parse additional user fields", "user", input.ID, "error", err)
		}
	}

	return input, nil
}
