package tmpl

import "embed"

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
)
