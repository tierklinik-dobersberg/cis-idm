package config

import "errors"

type Feature string

const (
	FeatureAll                 = "all"
	FeatureAddresses           = "addresses"
	FeatureEMails              = "emails"
	FeaturePhoneNumbers        = "phoneNumbers"
	FeatureEMailInvite         = "emailInvite"
	FeatureLoginByMail         = "loginByMail"
	FeatureAllowUsernameChange = "allowUsernameChange"
	FeatureSelfRegistration    = "registration"
)

var AllFeatures = []Feature{
	FeatureAddresses,
	FeatureEMails,
	FeaturePhoneNumbers,
	FeatureEMailInvite,
	FeatureLoginByMail,
	FeatureAllowUsernameChange,
	FeatureSelfRegistration,
}

var (
	ErrFeatureDisabled = errors.New("requested feature has been disabled")
)
