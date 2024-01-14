package selfservice

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/gofrs/uuid"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/pkg/data"
	"github.com/tierklinik-dobersberg/cis-idm/internal/bootstrap"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (svc *Service) GenerateAPIToken(ctx context.Context, req *connect.Request[idmv1.GenerateAPITokenRequest]) (*connect.Response[idmv1.GenerateAPITokenResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no jwt claims associated with request")
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return repo.RunInTransaction[*connect.Response[idmv1.GenerateAPITokenResponse]](ctx, svc.Datastore, func(tx *repo.Queries) (*connect.Response[idmv1.GenerateAPITokenResponse], error) {
		userRoles, err := tx.GetRolesForUser(ctx, claims.Subject)
		if err != nil {
			return nil, err
		}
		lm := data.IndexSlice(userRoles, func(r repo.Role) string {
			return r.ID
		})

		var roles = req.Msg.Roles
		if len(roles) > 0 {
			for _, r := range roles {
				if _, ok := lm[r]; !ok {
					return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("you're not allowed to use role %q", r))
				}
			}
		} else {
			// if no roles are specified we assign all user roles
			roles = maps.Keys(lm)
		}

		token, err := bootstrap.GenerateSecret(64)
		if err != nil {
			return nil, err
		}

		token = middleware.APITokenPrefix + token

		expiresAt := sql.NullTime{}

		if req.Msg.Expires.IsValid() {
			expiresAt = sql.NullTime{
				Time:  req.Msg.Expires.AsTime(),
				Valid: true,
			}
		}

		if err := tx.CreateAPIToken(ctx, repo.CreateAPITokenParams{
			ID:        id.String(),
			Token:     token,
			Name:      req.Msg.Description,
			UserID:    claims.Subject,
			ExpiresAt: expiresAt,
		}); err != nil {
			return nil, err
		}

		for _, roleId := range roles {
			if err := tx.AddRoleToToken(ctx, repo.AddRoleToTokenParams{
				TokenID: id.String(),
				RoleID:  roleId,
			}); err != nil {
				return nil, fmt.Errorf("failed to assign role %q to token: %w", roleId, err)
			}
		}

		return connect.NewResponse(&idmv1.GenerateAPITokenResponse{
			Token: &idmv1.APIToken{
				Id:            id.String(),
				Description:   req.Msg.Description,
				RedactedToken: token, // actually not redacted
				Expires:       req.Msg.Expires,
			},
		}), nil
	})
}

func (svc *Service) ListAPITokens(ctx context.Context, req *connect.Request[idmv1.ListAPITokensRequest]) (*connect.Response[idmv1.ListAPITokensResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no jwt claims associated with request context")
	}

	tokens, err := svc.Datastore.GetAPITokensForUser(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	res := &idmv1.ListAPITokensResponse{
		Tokens: make([]*idmv1.APIToken, len(tokens)),
	}

	for idx, token := range tokens {
		res.Tokens[idx] = &idmv1.APIToken{
			Id:            token.ID,
			Description:   token.Name,
			RedactedToken: token.Token[0:6],
			CreatedAt:     timestamppb.New(token.CreatedAt),
		}

		if token.ExpiresAt.Valid {
			res.Tokens[idx].Expires = timestamppb.New(token.ExpiresAt.Time)
		}
	}

	return connect.NewResponse(res), nil
}

func (svc *Service) RemoveAPIToken(ctx context.Context, req *connect.Request[idmv1.RemoveAPITokenRequest]) (*connect.Response[idmv1.RemoveAPITokenResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no jwt claims associated with request context")
	}

	res, err := svc.Datastore.RevokeUserAPIToken(ctx, repo.RevokeUserAPITokenParams{
		ID:     req.Msg.Id,
		UserID: claims.Subject,
	})

	if err != nil {
		return nil, err
	}

	if res == 0 {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("API token id not found"))
	}

	return connect.NewResponse(new(idmv1.RemoveAPITokenResponse)), nil
}
