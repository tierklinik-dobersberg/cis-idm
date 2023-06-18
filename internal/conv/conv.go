package conv

import (
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
)

type UserOption func(u *idmv1.Profile)

func UserProtoFromUser(user models.User) *idmv1.User {
	usr := &idmv1.User{
		Id:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Birthday:    user.Birthday,
	}

	return usr
}

func ProfileProtoFromUser(user models.User, useropts ...UserOption) *idmv1.Profile {
	profile := &idmv1.Profile{
		User: UserProtoFromUser(user),
	}

	for _, fn := range useropts {
		fn(profile)
	}

	return profile
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
		for _, nbr := range phoneNumbers {
			u.PhoneNumbers = append(u.PhoneNumbers, nbr.PhoneNumber)
		}
	}
}

func WithEmailAddresses(emails ...models.EMail) UserOption {
	return func(u *idmv1.Profile) {
		u.EmailAddresses = EmailProtosFromEmails(emails...)
	}
}
