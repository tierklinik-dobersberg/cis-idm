package roles

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/hashicorp/go-multierror"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

type Service struct {
	*app.Providers

	idmv1connect.UnimplementedRoleServiceHandler
}

func NewService(p *app.Providers) *Service {
	return &Service{
		Providers: p,
	}
}

func (svc *Service) CreateRole(ctx context.Context, req *connect.Request[idmv1.CreateRoleRequest]) (*connect.Response[idmv1.CreateRoleResponse], error) {
	roleModel := models.Role{
		Name:            req.Msg.Name,
		Description:     req.Msg.Description,
		DeleteProtected: req.Msg.DeleteProtection,
	}

	roleModel, err := svc.Datastore.CreateRole(ctx, roleModel)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.CreateRoleResponse{
		Role: conv.RoleProtoFromRole(roleModel),
	}), nil
}

func (svc *Service) UpdateRole(ctx context.Context, req *connect.Request[idmv1.UpdateRoleRequest]) (*connect.Response[idmv1.UpdateRoleResponse], error) {
	role, err := svc.Datastore.GetRoleByID(ctx, req.Msg.RoleId)
	if err != nil {
		if errors.Is(err, stmts.ErrNoResults) {
			return nil, connect.NewError(connect.CodeNotFound, nil)
		}

		return nil, err
	}

	paths := req.Msg.FieldMask.GetPaths()
	if len(paths) == 0 {
		paths = []string{"name", "description", "delete_protection"}
	}

	for _, p := range paths {
		switch p {
		case "name":
			role.Name = req.Msg.Name
		case "description":
			role.Description = req.Msg.Description
		case "delete_protection":
			role.DeleteProtected = req.Msg.DeleteProtection
		default:
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unknown field_mask.path %q", p))
		}
	}

	if err := svc.Datastore.UpdateRole(ctx, role); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.UpdateRoleResponse{
		Role: &idmv1.Role{
			Id:              role.ID,
			Description:     role.Description,
			DeleteProtected: role.DeleteProtected,
		},
	}), nil
}

func (svc *Service) DeleteRole(ctx context.Context, req *connect.Request[idmv1.DeleteRoleRequest]) (*connect.Response[idmv1.DeleteRoleResponse], error) {
	role, err := svc.Datastore.GetRoleByID(ctx, req.Msg.RoleId)
	if err != nil {
		if errors.Is(err, stmts.ErrNoResults) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}

		return nil, err
	}

	if role.DeleteProtected {
		return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("role is delete protected"))
	}

	if err := svc.Datastore.DeleteRole(ctx, role.ID); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.DeleteRoleResponse{}), nil
}

func (svc *Service) ListRoles(ctx context.Context, req *connect.Request[idmv1.ListRolesRequest]) (*connect.Response[idmv1.ListRolesResponse], error) {
	roles, err := svc.Datastore.GetRoles(ctx)
	if err != nil {
		return nil, err
	}

	res := &idmv1.ListRolesResponse{
		Roles: conv.RolesProtoFromRoles(roles...),
	}

	return connect.NewResponse(res), nil
}

func (svc *Service) GetRole(ctx context.Context, req *connect.Request[idmv1.GetRoleRequest]) (*connect.Response[idmv1.GetRoleResponse], error) {
	var (
		role models.Role
		err  error
	)

	selector := ""
	switch v := req.Msg.Search.(type) {
	case *idmv1.GetRoleRequest_Id:
		role, err = svc.Datastore.GetRoleByID(ctx, v.Id)
		selector = fmt.Sprintf("id=%q", v.Id)

	case *idmv1.GetRoleRequest_Name:
		role, err = svc.Datastore.GetRoleByName(ctx, v.Name)
		selector = fmt.Sprintf("name=%q", v.Name)

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	if err != nil {
		if errors.Is(err, stmts.ErrNoResults) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("role %s not found", selector))
		}
		return nil, err
	}

	return connect.NewResponse(&idmv1.GetRoleResponse{
		Role: conv.RoleProtoFromRole(role),
	}), nil
}

func (svc *Service) AssignRoleToUser(ctx context.Context, req *connect.Request[idmv1.AssignRoleToUserRequest]) (*connect.Response[idmv1.AssignRoleToUserResponse], error) {
	role, err := svc.Datastore.GetRoleByID(ctx, req.Msg.RoleId)
	if err != nil {
		if errors.Is(err, stmts.ErrNoResults) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("role not found"))
		}
		return nil, err
	}

	merr := new(multierror.Error)
	for _, userID := range req.Msg.UserId {
		user, err := svc.Datastore.GetUserByID(ctx, userID)
		if err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("user %s: %w", userID, err))
			continue
		}

		if err := svc.Datastore.AssignRoleToUser(ctx, user.ID, role.ID); err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("user %s: %w", userID, err))
			continue
		}
	}

	if err := merr.ErrorOrNil(); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.AssignRoleToUserResponse{}), nil
}

func (svc *Service) UnassignRoleFromUser(ctx context.Context, req *connect.Request[idmv1.UnassignRoleFromUserRequest]) (*connect.Response[idmv1.UnassignRoleFromUserResponse], error) {
	role, err := svc.Datastore.GetRoleByID(ctx, req.Msg.RoleId)
	if err != nil {
		if errors.Is(err, stmts.ErrNoResults) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("role not found"))
		}
		return nil, err
	}

	merr := new(multierror.Error)
	for _, userID := range req.Msg.UserId {
		user, err := svc.Datastore.GetUserByID(ctx, userID)
		if err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("user %s: %w", userID, err))
			continue
		}

		if err := svc.Datastore.UnassignRoleFromUser(ctx, user.ID, role.ID); err != nil {
			merr.Errors = append(merr.Errors, fmt.Errorf("user %s: %w", userID, err))
			continue
		}
	}

	if err := merr.ErrorOrNil(); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.UnassignRoleFromUserResponse{}), nil
}
