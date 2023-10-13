package conv

import (
	"context"
	"encoding/json"

	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"google.golang.org/protobuf/types/known/structpb"
)

type UserOption func(u *idmv1.Profile)

func UserProtoFromUser(ctx context.Context, user models.User) *idmv1.User {
	var extra *structpb.Struct
	if len(user.Extra) > 0 {
		var m map[string]any

		if err := json.Unmarshal([]byte(user.Extra), &m); err == nil {
			extra, err = structpb.NewStruct(m)
			if err != nil {
				log.L(ctx).Errorf("failed to encode user extra data: %s", err)
			}
		} else {
			log.L(ctx).Errorf("failed to decode user extra data: %s", err)
		}
	}

	usr := &idmv1.User{
		Id:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Birthday:    user.Birthday,
		Avatar:      user.Avatar,
		Extra:       extra,
	}

	return usr
}

func ProfileProtoFromUser(ctx context.Context, user models.User, useropts ...UserOption) *idmv1.Profile {
	profile := &idmv1.Profile{
		User:                UserProtoFromUser(ctx, user),
		TotpEnabled:         user.TOTPSecret != "",
		PasswordAuthEnabled: user.Password != "",
	}

	for _, fn := range useropts {
		fn(profile)
	}

	return profile
}

func RoleProtoFromRole(role models.Role) *idmv1.Role {
	return &idmv1.Role{
		Id:              role.ID,
		Name:            role.Name,
		Description:     role.Description,
		DeleteProtected: role.DeleteProtected,
	}
}

func RolesProtoFromRoles(roles ...models.Role) []*idmv1.Role {
	res := make([]*idmv1.Role, len(roles))

	for idx, r := range roles {
		res[idx] = RoleProtoFromRole(r)
	}

	return res
}

func EmailProtoFromEmail(email models.EMail) *idmv1.EMail {
	return &idmv1.EMail{
		Id:       email.ID,
		Address:  email.Address,
		Verified: email.Verified,
		Primary:  email.Primary,
	}
}

func EmailProtosFromEmails(emails ...models.EMail) []*idmv1.EMail {
	result := make([]*idmv1.EMail, len(emails))
	for idx, e := range emails {
		result[idx] = EmailProtoFromEmail(e)
	}

	return result
}

func AddressProtoFromAddress(addr models.Address) *idmv1.Address {
	return &idmv1.Address{
		Id:       addr.ID,
		CityCode: addr.CityCode,
		CityName: addr.CityName,
		Street:   addr.Street,
		Extra:    addr.Extra,
	}
}

func PhoneNumberProtoFromPhoneNumber(nbr models.PhoneNumber) *idmv1.PhoneNumber {
	return &idmv1.PhoneNumber{
		Id:       nbr.ID,
		Number:   nbr.PhoneNumber,
		Verified: nbr.Verified,
		Primary:  nbr.Primary,
	}
}

func PhoneNumberProtosFromPhoneNumbers(nbrs ...models.PhoneNumber) []*idmv1.PhoneNumber {
	result := make([]*idmv1.PhoneNumber, len(nbrs))
	for idx, n := range nbrs {
		result[idx] = PhoneNumberProtoFromPhoneNumber(n)
	}

	return result
}

func AddressProtosFromAddresses(addrs ...models.Address) []*idmv1.Address {
	result := make([]*idmv1.Address, len(addrs))
	for idx, a := range addrs {
		result[idx] = AddressProtoFromAddress(a)
	}

	return result
}

func WithAddresses(addresses ...models.Address) UserOption {
	return func(u *idmv1.Profile) {
		u.Addresses = AddressProtosFromAddresses(addresses...)
	}
}

func WithPhoneNumbers(phoneNumbers ...models.PhoneNumber) UserOption {
	return func(u *idmv1.Profile) {
		u.PhoneNumbers = PhoneNumberProtosFromPhoneNumbers(phoneNumbers...)
	}
}

func WithEmailAddresses(emails ...models.EMail) UserOption {
	return func(u *idmv1.Profile) {
		u.EmailAddresses = EmailProtosFromEmails(emails...)
	}
}

func WithPrimaryMail(mail *models.EMail) UserOption {
	return func(u *idmv1.Profile) {
		if mail == nil {
			return
		}

		u.User.PrimaryMail = EmailProtoFromEmail(*mail)
	}
}

func WithPrimaryPhone(phone *models.PhoneNumber) UserOption {
	return func(u *idmv1.Profile) {
		if phone == nil {
			return
		}

		u.User.PrimaryPhoneNumber = PhoneNumberProtoFromPhoneNumber(*phone)
	}
}

func WithRoles(roles ...models.Role) UserOption {
	return func(u *idmv1.Profile) {
		u.Roles = RolesProtoFromRoles(roles...)
	}
}

func WithUserHasRecoveryCodes(hasCodes bool) UserOption {
	return func(u *idmv1.Profile) {
		u.RecoveryCodesGenerated = hasCodes
	}
}
