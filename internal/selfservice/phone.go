package selfservice

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/bufbuild/connect-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/cis-idm/internal/cache"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
	"github.com/tierklinik-dobersberg/cis-idm/internal/sms"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
)

func (svc *Service) AddPhoneNumber(ctx context.Context, req *connect.Request[idmv1.AddPhoneNumberRequest]) (*connect.Response[idmv1.AddPhoneNumberResponse], error) {
	if !svc.Config.FeatureEnabled(config.FeaturePhoneNumbers) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("phone-numbers: %w", config.ErrFeatureDisabled))
	}

	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	m := models.PhoneNumber{
		UserID:      claims.Subject,
		PhoneNumber: req.Msg.Number,
		Verified:    false,
		Primary:     false,
	}

	m, err := svc.Datastore.AddUserPhoneNumber(ctx, m)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.AddPhoneNumberResponse{
		PhoneNumber: conv.PhoneNumberProtoFromPhoneNumber(m),
	}), nil
}

func (svc *Service) DeletePhoneNumber(ctx context.Context, req *connect.Request[idmv1.DeletePhoneNumberRequest]) (*connect.Response[idmv1.DeletePhoneNumberResponse], error) {
	if !svc.Config.FeatureEnabled(config.FeaturePhoneNumbers) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("phone-numbers: %w", config.ErrFeatureDisabled))
	}

	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	if err := svc.Datastore.DeleteUserPhoneNumber(ctx, claims.Subject, req.Msg.Id); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.DeletePhoneNumberResponse{}), nil
}

func (svc *Service) MarkPhoneNumberAsPrimary(ctx context.Context, req *connect.Request[idmv1.MarkPhoneNumberAsPrimaryRequest]) (*connect.Response[idmv1.MarkPhoneNumberAsPrimaryResponse], error) {
	if !svc.Config.FeatureEnabled(config.FeaturePhoneNumbers) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("phone-numbers: %w", config.ErrFeatureDisabled))
	}

	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	if err := svc.Datastore.MarkPhoneNumberAsPrimary(ctx, claims.Subject, req.Msg.Id); err != nil {
		return nil, err
	}

	return connect.NewResponse(&idmv1.MarkPhoneNumberAsPrimaryResponse{}), nil
}

func (svc *Service) ValidatePhoneNumber(ctx context.Context, req *connect.Request[idmv1.ValidatePhoneNumberRequest]) (*connect.Response[idmv1.ValidatePhoneNumberResponse], error) {
	if !svc.Config.FeatureEnabled(config.FeaturePhoneNumbers) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("phone-numbers: %w", config.ErrFeatureDisabled))
	}

	claims := middleware.ClaimsFromContext(ctx)
	if claims == nil {
		return nil, fmt.Errorf("no token claims associated with request context")
	}

	switch v := req.Msg.Step.(type) {
	case *idmv1.ValidatePhoneNumberRequest_Id:
		number, err := svc.Datastore.GetUserPhoneNumberByID(ctx, v.Id, claims.Subject)
		if err != nil {
			if errors.Is(err, stmts.ErrNoResults) {
				return nil, connect.NewError(connect.CodeNotFound, nil)
			}
		}

		if number.Verified {
			return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("number already verified"))
		}

		// Generate a new security code
		source := rand.NewSource(time.Now().UnixNano())
		rand := rand.New(source)
		code := fmt.Sprintf("%d", rand.Intn(999999-100000)+100000)

		// store the code in cache
		cacheKey := fmt.Sprintf("phone-verification:%s:%s", claims.Subject, code)
		if err := svc.Cache.PutKeyTTL(ctx, cacheKey, number.ID, time.Minute*5); err != nil {
			return nil, err
		}

		// Send a text message to the user
		if err := sms.SendTemplate(ctx, svc.SMSSender, svc.TemplateEngine, []string{number.PhoneNumber}, tmpl.VerifyPhoneNumber, &tmpl.VerifyPhoneNumberCtx{
			Code: code,
		}); err != nil {
			defer func() {
				_ = svc.Cache.DeleteKey(ctx, cacheKey)
			}()

			return nil, err
		}

		return connect.NewResponse(&idmv1.ValidatePhoneNumberResponse{}), nil

	case *idmv1.ValidatePhoneNumberRequest_Code:
		cacheKey := fmt.Sprintf("phone-verification:%s:%s", claims.Subject, v.Code)

		var numberID string
		if err := svc.Cache.GetAndDeleteKey(ctx, cacheKey, &numberID); err != nil {
			if errors.Is(err, cache.ErrKeyNotFound) {
				return nil, connect.NewError(connect.CodeAborted, fmt.Errorf("invalid security code"))
			}

			return nil, err
		}

		if err := svc.Datastore.MarkPhoneNumberAsVerified(ctx, claims.Subject, numberID); err != nil {
			if errors.Is(err, stmts.ErrNoResults) {
				return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("the verified number has been deleted from the profile"))
			}

			return nil, err
		}

		return connect.NewResponse(&idmv1.ValidatePhoneNumberResponse{}), nil

	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}
}
