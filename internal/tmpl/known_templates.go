package tmpl

import (
	"embed"

	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
)

//go:embed templates
var builtin embed.FS

type Kind string

const (
	KindSMS  = Kind("sms")
	KindMail = Kind("mail")
)

// Known is a known template that provides strong typing for
// the template variables and is bound by name.
type Known[C Context] struct {
	Kind Kind
	Name string
}

// Context types
type (
	VerifyPhoneNumberCtx struct {
		BaseContext
		Code string
	}

	SendPhoneSecurityCodeCtx struct {
		BaseContext
		Code string
	}

	RequestPasswordResetCtx struct {
		BaseContext
		User      models.User
		ResetLink string
	}

	VerifyMailCtx struct {
		BaseContext
		User       models.User
		VerifyLink string
	}

	InviteMailCtx struct {
		BaseContext
		RegisterURL string
		Name        string
		Inviter     models.User
	}
)

var (
	VerifyPhoneNumber = Known[*VerifyPhoneNumberCtx]{
		Name: "verify_phone_number",
		Kind: KindSMS,
	}

	SendPhoneSecurityCode = Known[*SendPhoneSecurityCodeCtx]{
		Name: "send_phone_security_code",
		Kind: KindSMS,
	}

	RequestPasswordReset = Known[*RequestPasswordResetCtx]{
		Name: "request_password_reset",
		Kind: KindMail,
	}

	VerifyMail = Known[*VerifyMailCtx]{
		Name: "verify_mail",
		Kind: KindMail,
	}

	InviteMail = Known[*InviteMailCtx]{
		Name: "user_invitation",
		Kind: KindMail,
	}
)
