package selfservice

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

func (svc *Service) GetRegisteredPasskeys(ctx context.Context, req *connect.Request[idmv1.GetRegisteredPasskeysRequest]) (*connect.Response[idmv1.GetRegisteredPasskeysResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	creds, err := svc.Datastore.GetWebauthnCreds(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	res := &idmv1.GetRegisteredPasskeysResponse{
		Passkeys: []*idmv1.RegisteredPasskey{},
	}

	for _, cred := range creds {
		res.Passkeys = append(res.Passkeys, &idmv1.RegisteredPasskey{
			Id:           cred.ID,
			ClientName:   cred.ClientName,
			ClientOs:     cred.ClientOs,
			ClientDevice: cred.ClientDevice,
			CredType:     cred.CredType,
		})
	}

	return connect.NewResponse(res), nil
}

func (svc *Service) RemovePasskey(ctx context.Context, req *connect.Request[idmv1.RemovePasskeyRequest]) (*connect.Response[idmv1.RemovePasskeyResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	rows, err := svc.Datastore.RemoveWebauthnCred(ctx, repo.RemoveWebauthnCredParams{UserID: claims.Subject, ID: req.Msg.Id})
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("passkey not found"))
	}

	return connect.NewResponse(&idmv1.RemovePasskeyResponse{}), nil
}
