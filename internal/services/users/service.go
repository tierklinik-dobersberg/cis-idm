package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ory/mail"
	"github.com/tidwall/sjson"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/bufbuild/connect-go"
	"github.com/hashicorp/go-multierror"
	"github.com/mennanov/fmutils"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/apis/pkg/data"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/common"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/mailer"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	idmv1connect.UnimplementedUserServiceHandler

	*app.Providers
}

func NewService(providers *app.Providers) (*Service, error) {
	svc := &Service{
		Providers: providers,
	}

	return svc, nil
}

func (svc *Service) ListUsers(ctx context.Context, req *connect.Request[idmv1.ListUsersRequest]) (*connect.Response[idmv1.ListUsersResponse], error) {
	users, err := svc.Datastore.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	res := &idmv1.ListUsersResponse{}
	for _, usr := range users {
		profileProto, err := svc.GetUserProfileProto(ctx, usr)
		if err != nil {
			log.L(ctx).Errorf("failed to get user profile for user %s: %s", usr.ID, err)
		}

		if len(req.Msg.FilterByRoles) > 0 {
			if !data.SliceOverlapsFunc(req.Msg.FilterByRoles, profileProto.Roles, func(role *idmv1.Role) string {
				return role.Id
			}) {
				continue
			}
		}

		res.Users = append(res.Users, profileProto)
	}

	// make sure we only include fields that are requested
	if fieldMaskPaths := req.Msg.GetFieldMask().GetPaths(); len(fieldMaskPaths) > 0 {
		if req.Msg.ExcludeFields {
			fmutils.Prune(res, fieldMaskPaths)
		} else {
			fmutils.Filter(res, fieldMaskPaths)
		}
	}

	return connect.NewResponse(res), nil
}

func (svc *Service) GetUser(ctx context.Context, req *connect.Request[idmv1.GetUserRequest]) (*connect.Response[idmv1.GetUserResponse], error) {
	var (
		user models.User
		err  error
	)

	switch v := req.Msg.Search.(type) {
	case *idmv1.GetUserRequest_Id:
		user, err = svc.Datastore.GetUserByID(ctx, v.Id)
	case *idmv1.GetUserRequest_Name:
		user, err = svc.Datastore.GetUserByName(ctx, v.Name)
	case *idmv1.GetUserRequest_Mail:
		user, _, err = svc.Datastore.GetUserByEMail(ctx, v.Mail)
	default:
		err = connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid request message"))
	}

	if err != nil {
		return nil, err
	}

	profileProto, err := svc.GetUserProfileProto(ctx, user)
	if err != nil {
		return nil, err
	}

	res := &idmv1.GetUserResponse{
		Profile: profileProto,
	}

	// make sure we only include fields that are requested
	if fieldMaskPaths := req.Msg.GetFieldMask().GetPaths(); len(fieldMaskPaths) > 0 {
		if req.Msg.ExcludeFields {
			fmutils.Prune(res, fieldMaskPaths)
		} else {
			fmutils.Filter(res, fieldMaskPaths)
		}
	}

	return connect.NewResponse(res), nil
}

func (svc *Service) DeleteUser(ctx context.Context, req *connect.Request[idmv1.DeleteUserRequest]) (*connect.Response[idmv1.DeleteUserResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("no jwt claims associated with request context"))
	}

	// make sure users just cannot directly delete their own profile.
	// TODO(ppacher): add support to "request profile deletion"
	if claims.Subject == req.Msg.Id {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("deleting your own account is not allowed"))
	}

	// actually delete the user from the repository
	if err := svc.Datastore.DeleteUser(ctx, req.Msg.Id); err != nil {
		return nil, err
	}

	// TODO(ppacher): invalidate all current access and refresh tokens

	return connect.NewResponse(&idmv1.DeleteUserResponse{}), nil
}

func (svc *Service) CreateUser(ctx context.Context, req *connect.Request[idmv1.CreateUserRequest]) (*connect.Response[idmv1.CreateUserResponse], error) {
	if req.Msg.Profile == nil || req.Msg.Profile.User == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid request: missing user object"))
	}
	usr := req.Msg.Profile.User

	userModel := models.User{
		Username:    usr.Username,
		FirstName:   usr.FirstName,
		LastName:    usr.LastName,
		DisplayName: usr.DisplayName,
		Avatar:      usr.Avatar,
		Birthday:    usr.Birthday,
	}

	if usr.Extra != nil {
		m := usr.Extra.AsMap()

		value, err := structpb.NewStruct(m)
		if err != nil {
			return nil, err
		}

		if err := svc.ValidateUserExtraData(value); err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}

		blob, err := json.Marshal(m)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to convert user extra data: %w", err))
		}

		userModel.Extra = string(blob)
	}

	if req.Msg.Password != "" {
		if req.Msg.PasswordIsBcrypt {
			userModel.Password = req.Msg.Password
		} else {
			hash, err := bcrypt.GenerateFromPassword([]byte(req.Msg.Password), bcrypt.DefaultCost)
			if err != nil {
				return nil, err
			}

			userModel.Password = string(hash)
		}
	}

	// actually create the user.
	userModel, err := svc.Datastore.CreateUser(ctx, userModel)
	if err != nil {
		return nil, err
	}

	merr := new(multierror.Error)

	// Create user addresses
	var userAddresses []models.Address
	if addresses := req.Msg.GetProfile().Addresses; len(addresses) > 0 {
		for _, addr := range addresses {
			addrModel := models.Address{
				UserID:   userModel.ID,
				CityCode: addr.CityCode,
				CityName: addr.CityName,
				Street:   addr.Street,
				Extra:    addr.Extra,
			}

			if addr, err := svc.Datastore.AddUserAddress(ctx, addrModel); err != nil {
				merr.Errors = append(merr.Errors, fmt.Errorf("failed to create user address: %w", err))
			} else {
				userAddresses = append(userAddresses, addr)
			}
		}
	}

	// Create phone number records
	var userPhoneNumbers []models.PhoneNumber
	if phoneNumbers := req.Msg.GetProfile().PhoneNumbers; len(phoneNumbers) > 0 {
		for _, nbr := range phoneNumbers {
			nbrModel := models.PhoneNumber{
				UserID:      userModel.ID,
				PhoneNumber: nbr.Number,
				Verified:    nbr.Verified,
				Primary:     nbr.Primary,
			}

			if phone, err := svc.Datastore.AddUserPhoneNumber(ctx, nbrModel); err != nil {
				merr.Errors = append(merr.Errors, fmt.Errorf("failed to create phone number: %w", err))
			} else {
				userPhoneNumbers = append(userPhoneNumbers, phone)
			}
		}
	}

	// create email address records
	var userEmails []models.EMail
	if emails := req.Msg.GetProfile().EmailAddresses; len(emails) > 0 {
		for idx, mail := range emails {
			mailModel := models.EMail{
				UserID:   userModel.ID,
				Address:  mail.Address,
				Verified: true,
				Primary:  idx == 0,
			}

			if email, err := svc.Datastore.CreateUserEmail(ctx, mailModel); err != nil {
				merr.Errors = append(merr.Errors, fmt.Errorf("failed to create email record: %w", err))
			} else {
				userEmails = append(userEmails, email)
			}
		}
	}

	// assign the user to all roles specified in the request.
	for _, role := range req.Msg.GetProfile().GetRoles() {
		err = errors.New("")

		if role.Id != "" {
			_, err = svc.Datastore.GetRoleByID(ctx, role.Id)
		}

		if err != nil || role.Id == "" {
			roleModel, err := svc.Datastore.GetRoleByName(ctx, role.Name)
			if err != nil {
				merr.Errors = append(merr.Errors, fmt.Errorf("role %q", role.Id))

				continue
			}

			role.Id = roleModel.ID
		}

		if err := svc.Datastore.AssignRoleToUser(ctx, userModel.ID, role.Id); err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("failed to assigne user %q to role %q", userModel.ID, role.Id))
		}
	}

	if err := merr.ErrorOrNil(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&idmv1.CreateUserResponse{
		Profile: conv.ProfileProtoFromUser(
			ctx,
			userModel,
			conv.WithAddresses(userAddresses...),
			conv.WithEmailAddresses(userEmails...),
			conv.WithPhoneNumbers(userPhoneNumbers...),
		),
	}), nil
}

func (svc *Service) UpdateUser(ctx context.Context, req *connect.Request[idmv1.UpdateUserRequest]) (*connect.Response[idmv1.UpdateUserResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no claims associated with request context")
	}

	user, err := svc.Datastore.GetUserByID(ctx, req.Msg.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user object: %w", err)
	}

	paths := req.Msg.GetFieldMask().GetPaths()
	if len(paths) == 0 {
		paths = []string{"username", "display_name", "first_name", "last_name", "avatar", "birthday"}
	}

	merr := new(multierror.Error)
	for _, p := range paths {
		switch p {
		case "username":
			if !svc.Config.FeatureEnabled(config.FeatureAllowUsernameChange) {
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("username changes are not allowed"))
			}

			user.Username = req.Msg.Username

			if len(user.Username) < 3 {
				merr.Errors = append(merr.Errors, fmt.Errorf("invalid username"))
			}
		case "display_name":
			user.DisplayName = req.Msg.DisplayName
			if len(user.DisplayName) < 3 {
				merr.Errors = append(merr.Errors, fmt.Errorf("invalid display-name"))
			}
		case "first_name":
			user.FirstName = req.Msg.FirstName
		case "last_name":
			user.LastName = req.Msg.LastName
		case "avatar":
			user.Avatar = req.Msg.Avatar
		case "birthday":
			user.Birthday = req.Msg.Birthday
		default:
			if strings.HasPrefix(p, "extra") {
				existing := make(map[string]any)

				if len(user.Extra) > 0 {
					if err := json.Unmarshal([]byte(user.Extra), &existing); err != nil {
						return nil, err
					}
				}

				updatedMap := req.Msg.Extra.AsMap()
				for key, value := range updatedMap {
					if value == nil {
						delete(existing, key)
						continue
					}

					existing[key] = value
				}

				value, err := structpb.NewStruct(existing)
				if err != nil {
					return nil, err
				}

				if err := svc.ValidateUserExtraData(value); err != nil {
					return nil, connect.NewError(connect.CodeInvalidArgument, err)
				}

				blob, err := json.Marshal(existing)
				if err != nil {
					return nil, err
				}

				user.Extra = string(blob)

				continue
			}

			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid field mask for update operation: invalid field name %q", p))
		}
	}

	if err := merr.ErrorOrNil(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := svc.Datastore.UpdateUser(ctx, user); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update user: %w", err))
	}

	profileProto, err := svc.GetUserProfileProto(ctx, user)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.UpdateUserResponse{
		Profile: profileProto,
	}), nil
}

func (svc *Service) InviteUser(ctx context.Context, req *connect.Request[idmv1.InviteUserRequest]) (*connect.Response[idmv1.InviteUserResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no claims associated with request context")
	}

	creator, err := svc.Datastore.GetUserByID(ctx, claims.Subject)
	if err != nil {
		if errors.Is(err, stmts.ErrNoResults) {
			return nil, fmt.Errorf("invalid jwt or account deleted")
		}

		return nil, err
	}

	common.EnsureDisplayName(&creator)

	var mailTemplates []*mail.Message

	merr := new(multierror.Error)
	for _, userInvite := range req.Msg.Invite {
		token, err := svc.Providers.GenerateRegistrationToken(ctx, creator, 1, 0, req.Msg.InitialRoles)
		if err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("failed to create token for %s: %w", userInvite.Email, err))
			continue
		}

		msg := mailer.Message{
			From: svc.Config.MailConfig.From,
			To:   []string{userInvite.Email},
		}

		if dr := svc.Providers.Config.DryRun; dr != nil && dr.MailTarget != "" {
			log.L(ctx).Infof("dry-run enabled, redirecting mails from %s to %s", userInvite.Email, dr.MailTarget)

			msg.To = []string{dr.MailTarget}
		}

		mail, err := mailer.PrepareTemplate(ctx, svc.Config, svc.TemplateEngine, msg, tmpl.InviteMail, &tmpl.InviteMailCtx{
			Name:        userInvite.Name,
			Inviter:     creator,
			RegisterURL: fmt.Sprintf(svc.Config.RegistrationURL, token, userInvite.Email, strings.ToLower(userInvite.Name)),
		})
		if err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("failed to send message to %s: %w", userInvite.Email, err))

			continue
		}

		mailTemplates = append(mailTemplates, mail)
	}

	if err := svc.Providers.Mailer.DialAndSend(mailTemplates...); err != nil {
		merr.Errors = append(merr.Errors, fmt.Errorf("failed to send mail: %w", err))
	}

	if err := merr.ErrorOrNil(); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.InviteUserResponse{}), nil
}

func (svc *Service) SetUserExtraKey(ctx context.Context, req *connect.Request[idmv1.SetUserExtraKeyRequest]) (*connect.Response[idmv1.SetUserExtraKeyResponse], error) {
	usr, err := svc.Datastore.GetUserByID(ctx, req.Msg.UserId)
	if err != nil {
		return nil, err
	}

	if usr.Extra == "" {
		usr.Extra = "{}"
	}

	if _, ok := req.Msg.Value.Kind.(*structpb.Value_NullValue); ok {
		usr.Extra, err = sjson.Delete(usr.Extra, req.Msg.Path)
	} else {
		usr.Extra, err = sjson.Set(usr.Extra, req.Msg.Path, req.Msg.Value.AsInterface())
	}

	if err != nil {
		return nil, err
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(usr.Extra), &m); err != nil {
		return nil, err
	}

	value, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}

	if err := svc.ValidateUserExtraData(value); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := svc.Datastore.UpdateUser(ctx, usr); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.SetUserExtraKeyResponse{}), nil
}

func (svc *Service) DeleteUserExtraKey(ctx context.Context, req *connect.Request[idmv1.DeleteUserExtraKeyRequest]) (*connect.Response[idmv1.DeleteUserExtraKeyResponse], error) {
	usr, err := svc.Datastore.GetUserByID(ctx, req.Msg.UserId)
	if err != nil {
		return nil, err
	}

	if usr.Extra != "" {
		usr.Extra, err = sjson.Delete(usr.Extra, req.Msg.Path)
		if err != nil {
			return nil, err
		}

		if err := svc.Datastore.UpdateUser(ctx, usr); err != nil {
			return nil, err
		}
	}

	return connect.NewResponse(&idmv1.DeleteUserExtraKeyResponse{}), nil
}
