package roles

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-multierror"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1/idmv1connect"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
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
	if !svc.Config.DynmicRolesEnabled() {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("dynamic role configuration is not enabled"))
	}

	if req.Msg.Id != "" {
		_, err := svc.Datastore.GetRoleByID(ctx, req.Msg.Id)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("failed to query for conflicting roles: %w", err)
			}
		}

		if err == nil {
			return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("role with id %q already exists", req.Msg.Id))
		}
	}

	params := repo.CreateRoleParams{
		ID:              req.Msg.Id,
		Name:            req.Msg.Name,
		Description:     req.Msg.Description,
		DeleteProtected: req.Msg.DeleteProtection,
	}

	if params.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}

		params.ID = id.String()
	}

	roleModel, err := svc.Datastore.CreateRole(ctx, params)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.CreateRoleResponse{
		Role: conv.RoleProtoFromRole(roleModel),
	}), nil
}

func (svc *Service) UpdateRole(ctx context.Context, req *connect.Request[idmv1.UpdateRoleRequest]) (*connect.Response[idmv1.UpdateRoleResponse], error) {
	if !svc.Config.DynmicRolesEnabled() {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("dynamic role configuration is not enabled"))
	}

	return repo.RunInTransaction(ctx, svc.Datastore, func(tx *repo.Queries) (*connect.Response[idmv1.UpdateRoleResponse], error) {
		role, err := tx.GetRoleByID(ctx, req.Msg.RoleId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, connect.NewError(connect.CodeNotFound, nil)
			}

			return nil, err
		}

		if role.Origin == "system" {
			return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("system roles cannot be modified"))
		}

		paths := req.Msg.FieldMask.GetPaths()
		if len(paths) == 0 {
			paths = []string{"name", "description", "delete_protection"}
		}

		update := repo.UpdateRoleParams{
			Name:            role.Name,
			Description:     role.Description,
			DeleteProtected: role.DeleteProtected,
			ID:              role.ID,
		}

		for _, p := range paths {
			switch p {
			case "name":
				update.Name = req.Msg.Name
			case "description":
				update.Description = req.Msg.Description
			case "delete_protection":
				update.DeleteProtected = req.Msg.DeleteProtection
			default:
				return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("unknown field_mask.path %q", p))
			}
		}

		role, err = tx.UpdateRole(ctx, update)
		if err != nil {
			return nil, err
		}

		return connect.NewResponse(&idmv1.UpdateRoleResponse{
			Role: &idmv1.Role{
				Id:              role.ID,
				Name:            role.Name,
				Description:     role.Description,
				DeleteProtected: role.DeleteProtected,
			},
		}), nil
	})
}

func (svc *Service) DeleteRole(ctx context.Context, req *connect.Request[idmv1.DeleteRoleRequest]) (*connect.Response[idmv1.DeleteRoleResponse], error) {
	if !svc.Config.DynmicRolesEnabled() {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("dynamic role management is not enabled"))
	}

	return repo.RunInTransaction(ctx, svc.Datastore, func(tx *repo.Queries) (*connect.Response[idmv1.DeleteRoleResponse], error) {
		role, err := tx.GetRoleByID(ctx, req.Msg.RoleId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, connect.NewError(connect.CodeNotFound, err)
			}

			return nil, err
		}

		if role.Origin == "system" {
			return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("system roles cannot be deleted"))
		}

		if role.DeleteProtected {
			return nil, connect.NewError(connect.CodeFailedPrecondition, fmt.Errorf("role is delete protected"))
		}

		rows, err := tx.DeleteRole(ctx, role.ID)
		if err != nil {
			return nil, err
		}

		if rows == 0 {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("role not found"))
		}

		return connect.NewResponse(&idmv1.DeleteRoleResponse{}), nil
	})
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
		role repo.Role
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("role %s not found", selector))
		}
		return nil, err
	}

	return connect.NewResponse(&idmv1.GetRoleResponse{
		Role: conv.RoleProtoFromRole(role),
	}), nil
}

func (svc *Service) AssignRoleToUser(ctx context.Context, req *connect.Request[idmv1.AssignRoleToUserRequest]) (*connect.Response[idmv1.AssignRoleToUserResponse], error) {
	return repo.RunInTransaction(ctx, svc.Datastore, func(tx *repo.Queries) (*connect.Response[idmv1.AssignRoleToUserResponse], error) {
		role, err := tx.GetRoleByID(ctx, req.Msg.RoleId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("role not found"))
			}
			return nil, err
		}

		merr := new(multierror.Error)
		for _, userID := range req.Msg.UserId {
			user, err := tx.GetUserByID(ctx, userID)
			if err != nil {
				merr.Errors = append(merr.Errors, fmt.Errorf("user %s: %w", userID, err))
				continue
			}

			if err := tx.AssignRoleToUser(ctx, repo.AssignRoleToUserParams{
				UserID: user.ID,
				RoleID: role.ID,
			}); err != nil {
				merr.Errors = append(merr.Errors, fmt.Errorf("user %s: %w", userID, err))
				continue
			}
		}

		if err := merr.ErrorOrNil(); err != nil {
			return nil, err
		}

		return connect.NewResponse(&idmv1.AssignRoleToUserResponse{}), nil
	})
}

func (svc *Service) UnassignRoleFromUser(ctx context.Context, req *connect.Request[idmv1.UnassignRoleFromUserRequest]) (*connect.Response[idmv1.UnassignRoleFromUserResponse], error) {
	return repo.RunInTransaction(ctx, svc.Datastore, func(tx *repo.Queries) (*connect.Response[idmv1.UnassignRoleFromUserResponse], error) {
		role, err := tx.GetRoleByID(ctx, req.Msg.RoleId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("role not found"))
			}
			return nil, err
		}

		merr := new(multierror.Error)
		for _, userID := range req.Msg.UserId {
			user, err := tx.GetUserByID(ctx, userID)
			if err != nil {
				merr.Errors = append(merr.Errors, fmt.Errorf("user %s: %w", userID, err))
				continue
			}

			rows, err := tx.UnassignRoleFromUser(ctx, repo.UnassignRoleFromUserParams{
				UserID: user.ID,
				RoleID: role.ID,
			})
			if err != nil {
				merr.Errors = append(merr.Errors, fmt.Errorf("user %s: %w", userID, err))
				continue
			}

			if rows == 0 {
				merr.Errors = append(merr.Errors, fmt.Errorf("user-assignment for user %s and role %s: not found", user.ID, role.ID))
			}
		}

		if err := merr.ErrorOrNil(); err != nil {
			return nil, err
		}

		return connect.NewResponse(&idmv1.UnassignRoleFromUserResponse{}), nil
	})
}

func (svc *Service) ResolveRolePermissions(ctx context.Context, req *connect.Request[idmv1.ResolveRolePermissionsRequest]) (*connect.Response[idmv1.ResolveRolePermissionsResponse], error) {
	return repo.RunInTransaction(ctx, svc.Datastore, func(tx *repo.Queries) (*connect.Response[idmv1.ResolveRolePermissionsResponse], error) {
		role, err := tx.GetRoleByID(ctx, req.Msg.RoleId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("role id not fuond"))
			}

			return nil, err
		}

		permissions, err := tx.GetRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get role permissions: %w", err)
		}

		resolved, err := svc.Config.Permissions.Resolve(permissions)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve role permissions: %w", err)
		}

		return connect.NewResponse(&idmv1.ResolveRolePermissionsResponse{
			Permissions: resolved,
		}), nil
	}, repo.ReadOnly())
}
