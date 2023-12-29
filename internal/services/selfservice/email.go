package selfservice

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

func (svc *Service) AddEmailAddress(ctx context.Context, req *connect.Request[idmv1.AddEmailAddressRequest]) (*connect.Response[idmv1.AddEmailAddressResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	user, err := svc.Datastore.GetUserByID(ctx, claims.Subject)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
	}

	addedMail, mails, err := svc.Common.AddEmailAddressToUser(ctx, repo.UserEmail{
		UserID:  claims.Subject,
		Address: req.Msg.Email,
	})
	if err != nil {
		return nil, err
	}

	if err := svc.SendMailVerification(ctx, user, *addedMail); err != nil {
		return nil, err
	}

	res := connect.NewResponse(&idmv1.AddEmailAddressResponse{
		Emails: conv.EmailProtosFromEmails(mails...),
	})

	return res, nil
}

func (svc *Service) DeleteEmailAddress(ctx context.Context, req *connect.Request[idmv1.DeleteEmailAddressRequest]) (*connect.Response[idmv1.DeleteEmailAddressResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	mails, err := svc.Common.DeleteEmailAddressFromUser(ctx, claims.Subject, req.Msg.Id)
	if err != nil {
		return nil, err
	}

	res := connect.NewResponse(&idmv1.DeleteEmailAddressResponse{
		Emails: conv.EmailProtosFromEmails(mails...),
	})

	return res, nil
}

func (svc *Service) MarkEmailAsPrimary(ctx context.Context, req *connect.Request[idmv1.MarkEmailAsPrimaryRequest]) (*connect.Response[idmv1.MarkEmailAsPrimaryResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	if err := svc.Common.MarkEmailAsPrimary(ctx, claims.Subject, req.Msg.Id); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.MarkEmailAsPrimaryResponse{}), nil
}

func (svc *Service) ValidateEmail(ctx context.Context, req *connect.Request[idmv1.ValidateEmailRequest]) (*connect.Response[idmv1.ValidateEmailResponse], error) {
	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	switch v := req.Msg.Kind.(type) {
	case *idmv1.ValidateEmailRequest_EmailId:
		email, err := svc.Datastore.GetEmailByID(ctx, repo.GetEmailByIDParams{
			UserID: claims.Subject,
			ID:     v.EmailId,
		})
		if err != nil {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("unknown email id"))
		}

		user, err := svc.Datastore.GetUserByID(ctx, claims.Subject)
		if err != nil {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("unknown user"))
		}

		if email.Verified {
			return nil, connect.NewError(connect.CodeAborted, fmt.Errorf("email already verified"))
		}

		if err := svc.SendMailVerification(ctx, user, email); err != nil {
			return nil, err
		}

	case *idmv1.ValidateEmailRequest_Token:
		cacheKey := fmt.Sprintf("verify-email:%s:%s", claims.Subject, v.Token)
		var emailID string
		if err := svc.Cache.GetAndDeleteKey(ctx, cacheKey, &emailID); err != nil {
			return nil, err
		}

		email, err := svc.Datastore.GetEmailByID(ctx, repo.GetEmailByIDParams{
			UserID: claims.Subject,
			ID:     emailID,
		})
		if err != nil {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("unknown email id"))
		}

		rows, err := svc.Datastore.MarkEmailAsVerified(ctx, repo.MarkEmailAsVerifiedParams{UserID: claims.Subject, ID: email.ID})
		if err != nil {
			return nil, err
		}

		if rows == 0 {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user mail not found"))
		}

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid request"))
	}

	return connect.NewResponse(&idmv1.ValidateEmailResponse{}), nil
}
