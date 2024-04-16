package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
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
	"github.com/tierklinik-dobersberg/cis-idm/internal/mailer"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
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

func (svc *Service) Impersonate(ctx context.Context, req *connect.Request[idmv1.ImpersonateRequest]) (*connect.Response[idmv1.ImpersonateResponse], error) {
	user, err := svc.Datastore.GetUserByID(ctx, req.Msg.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user %q not found", req.Msg.UserId))
		}
	}

	roles, err := svc.Datastore.GetRolesForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	tokenMessage := &idmv1.ImpersonateResponse{}
	res := connect.NewResponse(tokenMessage)

	token, _, err := svc.AddAccessToken(user, roles, svc.Config.AccessTTL(), "", "impersonate", res.Header())
	if err != nil {
		return nil, err
	}

	authUser := middleware.ClaimsFromContext(ctx)

	log.L(ctx).Infof("user %s (%q) impersonated %s (%q)", authUser.Subject, authUser.Name, user.ID, user.Username)

	tokenMessage.AccessToken = token

	return res, nil
}

func (svc *Service) ListUsers(ctx context.Context, req *connect.Request[idmv1.ListUsersRequest]) (*connect.Response[idmv1.ListUsersResponse], error) {
	users, err := svc.Datastore.GetAllUsers(ctx)
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
			if !data.ElemInBothSlicesFunc(req.Msg.FilterByRoles, profileProto.Roles, func(role *idmv1.Role) string {
				return role.Id
			}) {
				continue
			}
		}

		res.Users = append(res.Users, profileProto)
	}

	// make sure we only include fields that are requested
	if fieldMaskPaths := req.Msg.GetFieldMask().GetPaths(); len(fieldMaskPaths) > 0 {
		// compatibility: unfortunatley the field in ListUsersResponse is called
		// users instead of profiles (which would be better for consistency).
		// It seems this is a common error cause since users expect the field to be
		// called "profiles".
		// Make sure we accept both versions for the field mask.
		for idx, path := range fieldMaskPaths {
			if strings.HasPrefix(path, "profiles") {
				fieldMaskPaths[idx] = strings.Replace(path, "profiles", "users", 1)
			}
		}

		if req.Msg.ExcludeFields {
			fmutils.Prune(res, fieldMaskPaths)
		} else {
			fmutils.Filter(res, fieldMaskPaths)
		}
	}

	return connect.NewResponse(res), nil
}

func (svc *Service) SetUserPassword(ctx context.Context, req *connect.Request[idmv1.SetUserPasswordRequest]) (*connect.Response[idmv1.SetUserPasswordResponse], error) {
	_, err := svc.Datastore.GetUserByID(ctx, req.Msg.UserId)
	if err != nil {
		return nil, err
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Msg.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}

	rows, err := svc.Datastore.SetUserPassword(ctx, repo.SetUserPasswordParams{ID: req.Msg.UserId, Password: string(newHashedPassword)})
	if err != nil {
		return nil, fmt.Errorf("failed to save user password: %w", err)
	}

	if rows == 0 {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
	}

	return connect.NewResponse(new(idmv1.SetUserPasswordResponse)), nil
}

func (svc *Service) GetUser(ctx context.Context, req *connect.Request[idmv1.GetUserRequest]) (*connect.Response[idmv1.GetUserResponse], error) {
	var (
		user repo.User
		err  error
	)

	switch v := req.Msg.Search.(type) {
	case *idmv1.GetUserRequest_Id:
		user, err = svc.Datastore.GetUserByID(ctx, v.Id)
	case *idmv1.GetUserRequest_Name:
		user, err = svc.Datastore.GetUserByName(ctx, v.Name)
	case *idmv1.GetUserRequest_Mail:
		var response repo.GetUserByEMailRow
		response, err = svc.Datastore.GetUserByEMail(ctx, v.Mail)

		user = response.User
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
	rows, err := svc.Datastore.DeleteUser(ctx, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
	}

	// TODO(ppacher): invalidate all current access and refresh tokens

	return connect.NewResponse(&idmv1.DeleteUserResponse{}), nil
}

func (svc *Service) CreateUser(ctx context.Context, req *connect.Request[idmv1.CreateUserRequest]) (*connect.Response[idmv1.CreateUserResponse], error) {
	if req.Msg.Profile == nil || req.Msg.Profile.User == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid request: missing user object"))
	}
	usr := req.Msg.Profile.User

	userModel := repo.User{
		ID:          usr.Id,
		Username:    usr.Username,
		FirstName:   usr.FirstName,
		LastName:    usr.LastName,
		DisplayName: usr.DisplayName,
		Avatar:      usr.Avatar,
		Birthday:    usr.Birthday,
	}

	if userModel.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}

		userModel.ID = id.String()
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

	return repo.RunInTransaction(ctx, svc.Datastore, func(tx *repo.Queries) (*connect.Response[idmv1.CreateUserResponse], error) {

		// actually create the user.
		userModel, err := tx.CreateUser(ctx, repo.CreateUserParams{
			ID:          userModel.ID,
			Username:    userModel.Username,
			DisplayName: userModel.DisplayName,
			FirstName:   userModel.FirstName,
			LastName:    userModel.LastName,
			Extra:       userModel.Extra,
			Avatar:      userModel.Avatar,
			Birthday:    userModel.Birthday,
			Password:    userModel.Password,
		})
		if err != nil {
			return nil, err
		}

		merr := new(multierror.Error)

		// Create user addresses
		if addresses := req.Msg.GetProfile().Addresses; len(addresses) > 0 {
			for _, addr := range addresses {
				id, err := uuid.NewV4()
				if err != nil {
					return nil, err
				}

				addrModel := repo.CreateUserAddressParams{
					ID:       id.String(),
					UserID:   userModel.ID,
					CityCode: addr.CityCode,
					CityName: addr.CityName,
					Street:   addr.Street,
					Extra:    addr.Extra,
				}

				if _, err := tx.CreateUserAddress(ctx, addrModel); err != nil {
					merr.Errors = append(merr.Errors, fmt.Errorf("failed to create user address: %w", err))
				}
			}
		}

		// Create phone number records
		if phoneNumbers := req.Msg.GetProfile().PhoneNumbers; len(phoneNumbers) > 0 {
			for _, nbr := range phoneNumbers {
				id, err := uuid.NewV4()
				if err != nil {
					return nil, err
				}

				nbrModel := repo.CreateUserPhoneNumberParams{
					ID:          id.String(),
					UserID:      userModel.ID,
					PhoneNumber: nbr.Number,
					Verified:    nbr.Verified,
				}

				if _, err := tx.CreateUserPhoneNumber(ctx, nbrModel); err != nil {
					merr.Errors = append(merr.Errors, fmt.Errorf("failed to create phone number: %w", err))
				}
			}
		}

		// create email address records
		if emails := req.Msg.GetProfile().EmailAddresses; len(emails) > 0 {
			for idx, mail := range emails {
				id, err := uuid.NewV4()
				if err != nil {
					return nil, err
				}

				mailModel := repo.CreateEMailParams{
					ID:        id.String(),
					UserID:    userModel.ID,
					Address:   mail.Address,
					Verified:  true,
					IsPrimary: idx == 0,
				}

				if _, err := tx.CreateEMail(ctx, mailModel); err != nil {
					merr.Errors = append(merr.Errors, fmt.Errorf("failed to create email record: %w", err))
				}
			}
		}

		// assign the user to all roles specified in the request.
		for _, role := range req.Msg.GetProfile().GetRoles() {
			err = errors.New("")

			if role.Id != "" {
				_, err = tx.GetRoleByID(ctx, role.Id)
			}

			if err != nil || role.Id == "" {
				roleModel, err := tx.GetRoleByName(ctx, role.Name)
				if err != nil {
					merr.Errors = append(merr.Errors, fmt.Errorf("role %q", role.Id))

					continue
				}

				role.Id = roleModel.ID
			}

			if err := tx.AssignRoleToUser(ctx, repo.AssignRoleToUserParams{UserID: userModel.ID, RoleID: role.Id}); err != nil {
				merr.Errors = append(merr.Errors, fmt.Errorf("failed to assigne user %q to role %q", userModel.ID, role.Id))
			}
		}

		if err := merr.ErrorOrNil(); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		profile, err := app.GetUserProfileProto(ctx, tx, svc.Config, userModel)
		if err != nil {
			return nil, err
		}

		return connect.NewResponse(&idmv1.CreateUserResponse{
			Profile: profile,
		}), nil
	})
}

func (svc *Service) UpdateUser(ctx context.Context, req *connect.Request[idmv1.UpdateUserRequest]) (*connect.Response[idmv1.UpdateUserResponse], error) {
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
			if !svc.Config.AllowUsernameChange {
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

	if _, err := svc.Datastore.UpdateUser(ctx, repo.UpdateUserParams{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Extra:       user.Extra,
		Avatar:      user.Avatar,
		Birthday:    user.Birthday,
		ID:          user.ID,
	}); err != nil {
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
		if errors.Is(err, sql.ErrNoRows) {
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
			RegisterURL: fmt.Sprintf(svc.Config.UserInterface.RegistrationURL, token, userInvite.Email, strings.ToLower(userInvite.Name)),
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

	rows, err := svc.Datastore.SetUserExtraData(ctx, repo.SetUserExtraDataParams{
		Extra: usr.Extra,
		ID:    usr.ID,
	})
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
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

		rows, err := svc.Datastore.SetUserExtraData(ctx, repo.SetUserExtraDataParams{
			Extra: "",
			ID:    usr.ID,
		})
		if err != nil {
			return nil, err
		}

		if rows == 0 {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
		}
	}

	return connect.NewResponse(&idmv1.DeleteUserExtraKeyResponse{}), nil
}

func (svc *Service) SendAccountCreationNotice(ctx context.Context, req *connect.Request[idmv1.SendAccountCreationNoticeRequest]) (*connect.Response[idmv1.SendAccountCreationNoticeResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no request claims associated with request")
	}

	creator, err := svc.Providers.Datastore.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("failed to find authenticated user in database"))
	}

	merr := new(multierror.Error)

	for _, userId := range req.Msg.UserIds {
		target, err := svc.Datastore.GetUserByID(ctx, userId)
		if err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("failed to find user %q: %w", userId, err))

			continue
		}

		primaryMail, err := svc.Datastore.GetPrimaryEmailForUserByID(ctx, userId)
		if err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("failed to get primary email for user %q: %w", userId, err))

			continue
		}

		code, cacheKey, err := svc.Common.GeneratePasswordResetToken(ctx, target.ID)
		if err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("failed to generate password reset token for user %q: %w", userId, err))

			continue
		}

		// Send a text message to the user
		msg := mailer.Message{
			From: svc.Config.MailConfig.From,
			To:   []string{primaryMail.Address},
		}

		if err := mailer.SendTemplate(ctx, svc.Config, svc.TemplateEngine, svc.Mailer, msg, tmpl.AccountCreationNotice, &tmpl.AccountCreationNoticeCtx{
			Creator:   creator,
			User:      target,
			ResetLink: fmt.Sprintf(svc.Config.UserInterface.PasswordResetURL, code),
		}); err != nil {
			defer func() {
				_ = svc.Cache.DeleteKey(ctx, cacheKey)
			}()

			merr.Errors = append(merr.Errors, fmt.Errorf("failed to send account creation notice to user %q: %w", userId, err))

			log.L(ctx).Errorf("failed to send account creation notice to user %q: %s", userId, err)
		}
	}

	if err := merr.ErrorOrNil(); err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	return connect.NewResponse(new(idmv1.SendAccountCreationNoticeResponse)), nil
}

func (svc *Service) ResolveUserPermissions(ctx context.Context, req *connect.Request[idmv1.ResolveUserPermissionsRequest]) (*connect.Response[idmv1.ResolveUserPermissionsResponse], error) {
	return repo.RunInTransaction(ctx, svc.Datastore, func(tx *repo.Queries) (*connect.Response[idmv1.ResolveUserPermissionsResponse], error) {
		user, err := tx.GetUserByID(ctx, req.Msg.UserId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user id not fuond"))
			}

			return nil, err
		}

		roles, err := tx.GetRolesForUser(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user roles: %w", err)
		}

		var permissions []string
		for _, r := range roles {
			rolePerms, err := tx.GetRolePermissions(ctx, r.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get role permissions: %w", err)
			}

			permissions = append(permissions, rolePerms...)
		}

		resolved, err := svc.Config.PermissionTree().Resolve(permissions)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve role permissions: %w", err)
		}

		return connect.NewResponse(&idmv1.ResolveUserPermissionsResponse{
			Permissions: resolved,
		}), nil
	}, repo.ReadOnly())
}
