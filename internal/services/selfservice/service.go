package selfservice

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/hashicorp/go-multierror"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

type Service struct {
	idmv1connect.UnimplementedSelfServiceServiceHandler

	*app.Providers
}

func NewService(provider *app.Providers) *Service {
	return &Service{
		Providers: provider,
	}
}

func (svc *Service) UpdateProfile(ctx context.Context, req *connect.Request[idmv1.UpdateProfileRequest]) (*connect.Response[idmv1.UpdateProfileResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no claims associated with request context")
	}

	user, err := svc.Datastore.GetUserByID(ctx, claims.Subject)
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
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid field mask for update operation: invalid field name %q", p))
		}
	}

	if err := merr.ErrorOrNil(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	user, err = svc.Datastore.UpdateUser(ctx, repo.UpdateUserParams{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Extra:       user.Extra,
		Avatar:      user.Avatar,
		Birthday:    user.Birthday,
		ID:          user.ID,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update user: %w", err))
	}

	return connect.NewResponse(&idmv1.UpdateProfileResponse{
		User: conv.UserProtoFromUser(ctx, user),
	}), nil
}
