// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package repo

import (
	"database/sql"
	"time"
)

type MfaBackupCode struct {
	Code   string
	UserID string
}

type Permission struct {
	Name        string
	Description sql.NullString
}

type RegistrationToken struct {
	Token        string
	Expires      sql.NullTime
	AllowedUsage sql.NullInt64
	InitialRoles string
	CreatedBy    string
	CreatedAt    time.Time
}

type Role struct {
	ID              string
	Name            string
	Description     string
	DeleteProtected bool
}

type RoleAssignment struct {
	UserID string
	RoleID string
}

type RolePermission struct {
	Permission string
	RoleID     string
}

type TokenInvalidation struct {
	TokenID   string
	UserID    string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type User struct {
	ID          string
	Username    string
	DisplayName string
	FirstName   string
	LastName    string
	Extra       string
	Avatar      string
	Birthday    string
	Password    string
	TotpSecret  sql.NullString
}

type UserAddress struct {
	ID       string
	UserID   string
	CityCode string
	CityName string
	Street   string
	Extra    string
}

type UserApiToken struct {
	Token     string
	Name      string
	UserID    string
	ExpiresAt sql.NullTime
	CreatedAt time.Time
}

type UserApiTokenRole struct {
	Token  string
	RoleID string
}

type UserEmail struct {
	ID        string
	UserID    string
	Address   string
	Verified  bool
	IsPrimary bool
}

type UserPhoneNumber struct {
	ID          string
	UserID      string
	PhoneNumber string
	IsPrimary   bool
	Verified    bool
}

type WebauthnCred struct {
	ID           string
	UserID       string
	Cred         string
	CredType     string
	ClientName   string
	ClientOs     string
	ClientDevice string
}

type WebpushSubscription struct {
	ID        string
	UserID    string
	UserAgent string
	Endpoint  string
	Auth      string
	Key       string
	TokenID   string
}
