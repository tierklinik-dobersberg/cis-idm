package conv

import (
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
)

type UserOption func(u *idmv1.User)

func UserProtoFromUser(user models.User, useropts ...UserOption) *idmv1.User {
	usr := &idmv1.User{
		Id:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
	}

	for _, fn := range useropts {
		fn(usr)
	}

	return usr
}

func WithAddresses(addresses ...models.Address) UserOption {
	return func(u *idmv1.User) {
		for _, addr := range addresses {
			u.Addresses = append(u.Addresses, &idmv1.Address{
				CityCode: addr.CityCode,
				CityName: addr.CityName,
				Street:   addr.Street,
				Extra:    addr.Extra,
			})
		}
	}
}

func WithPhoneNumbers(phoneNumbers ...models.PhoneNumber) UserOption {
	return func(u *idmv1.User) {
		for _, nbr := range phoneNumbers {
			u.PhoneNumbers = append(u.PhoneNumbers, nbr.PhoneNumber)
		}
	}
}

func WithEmailAddresses(emails ...models.EMail) UserOption {
	return func(u *idmv1.User) {
		for _, mail := range emails {
			u.EmailAddresses = append(u.EmailAddresses, &idmv1.EMail{
				Address:  mail.Address,
				Verified: mail.Verified,
			})
		}
	}
}
