package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/bufbuild/protovalidate-go"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/bootstrap"
	"github.com/tierklinik-dobersberg/cis-idm/internal/cache"
	"github.com/tierklinik-dobersberg/cis-idm/internal/common"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/conv"
	"github.com/tierklinik-dobersberg/cis-idm/internal/mailer"
	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/sms"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/structpb"
)

type Providers struct {
	TemplateEngine *tmpl.Engine
	SMSSender      sms.Sender
	Mailer         mailer.Mailer
	Datastore      *repo.Queries
	Config         config.Config
	Common         *common.Service
	ProtoRegistry  *protoregistry.Files
	Validator      *protovalidate.Validator
	Cache          cache.Cache
}

func (p *Providers) SendMailVerification(ctx context.Context, user repo.User, mail repo.UserEmail) error {
	secret, err := bootstrap.GenerateSecret(16)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("verify-email:%s:%s", user.ID, secret)
	if err := p.Cache.PutKeyTTL(ctx, cacheKey, mail.ID, time.Hour*24); err != nil {
		return err
	}

	common.EnsureDisplayName(&user)

	msg := mailer.Message{
		From: p.Config.MailConfig.From,
		To:   []string{mail.Address},
	}
	if err := mailer.SendTemplate(ctx, p.Config, p.TemplateEngine, p.Mailer, msg, tmpl.VerifyMail, &tmpl.VerifyMailCtx{
		User:       user,
		VerifyLink: fmt.Sprintf(p.Config.VerifyMailURL, secret),
	}); err != nil {
		return err
	}

	return nil
}

func (p *Providers) GenerateRegistrationToken(ctx context.Context, creator repo.User, maxCount uint64, ttl time.Duration, initialRoles []string) (string, error) {
	token, err := bootstrap.GenerateSecret(8)
	if err != nil {
		return "", err
	}

	tokenModel := repo.CreateRegistrationTokenParams{
		Token:     token,
		CreatedBy: creator.ID,
		CreatedAt: time.Now(),
	}

	if maxCount > 0 {
		tokenModel.AllowedUsage = sql.NullInt64{
			Int64: int64(maxCount),
			Valid: true,
		}
	}

	if ttl > 0 {
		expires := time.Now().Add(ttl)
		tokenModel.Expires = sql.NullTime{
			Time:  expires,
			Valid: true,
		}
	}

	if len(initialRoles) > 0 {
		var initialRoleIDs []string

		for _, role := range initialRoles {
			roleModel, err := p.Datastore.GetRoleByID(ctx, role)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					roleModel, err = p.Datastore.GetRoleByName(ctx, role)
				}
			}

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return "", connect.NewError(connect.CodeNotFound, fmt.Errorf("roles %q does not exist", role))
				}

				return "", err
			}
			initialRoleIDs = append(initialRoleIDs, roleModel.ID)
		}

		roleBlob, err := json.Marshal(initialRoleIDs)
		if err != nil {
			return "", err
		}

		tokenModel.InitialRoles = string(roleBlob)
	}

	if err := p.Datastore.CreateRegistrationToken(ctx, tokenModel); err != nil {
		return "", err
	}

	return tokenModel.Token, nil
}

func getCurrentFieldVisiblity(ctx context.Context, id string) config.FieldVisibility {
	if claims := middleware.ClaimsFromContext(ctx); claims != nil {
		if claims.AppMetadata != nil && claims.AppMetadata.Authorization != nil && slices.Contains(claims.AppMetadata.Authorization.Roles, "idm_superuser") {
			return config.FieldVisibilityPrivate
		} else if claims.Subject == id {
			return config.FieldVisibilitySelf
		} else {
			return config.FieldVisibilityAuthenticated
		}
	}

	return config.FieldVisibilityPublic
}

func (p *Providers) GetUserProfileProto(ctx context.Context, usr repo.User) (*idmv1.Profile, error) {
	addresses, err := p.Datastore.GetUserAddresses(ctx, usr.ID)
	if err != nil {
		log.L(ctx).Errorf("failed to get user addresses: %s", err)
	}

	phones, err := p.Datastore.GetPhoneNumbersByUserID(ctx, usr.ID)
	if err != nil {
		log.L(ctx).Errorf("failed to get user phone numbers: %s", err)
	}

	roles, err := p.Datastore.GetRolesForUser(ctx, usr.ID)
	if err != nil {
		log.L(ctx).Errorf("failed to get user roles: %s", err)
	}

	emails, err := p.Datastore.GetEmailsForUserByID(ctx, usr.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load emails: %w", err)
	}

	hasBackupCodes, err := p.Datastore.UserHasTOTPEnrolled(ctx, usr.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing recovery codes: %w", err)
	}

	var primaryMail *repo.UserEmail
	for _, mail := range emails {
		if mail.IsPrimary {
			primaryMail = new(repo.UserEmail)
			*primaryMail = mail

			break
		}
	}

	var primaryPhone *repo.UserPhoneNumber
	for _, phone := range phones {
		if phone.IsPrimary {
			primaryPhone = new(repo.UserPhoneNumber)
			*primaryPhone = phone

			break
		}
	}

	profile := conv.ProfileProtoFromUser(
		ctx,
		usr,
		conv.WithAddresses(addresses...),
		conv.WithEmailAddresses(emails...),
		conv.WithPhoneNumbers(phones...),
		conv.WithRoles(roles...),
		conv.WithPrimaryMail(primaryMail),
		conv.WithPrimaryPhone(primaryPhone),
		conv.WithUserHasRecoveryCodes(hasBackupCodes),
	)

	if extra := profile.GetUser().GetExtra(); extra != nil {
		currentVisiblity := getCurrentFieldVisiblity(ctx, usr.ID)

		for key, propertyConfig := range p.Config.ExtraDataConfig {
			value := extra.Fields[key]
			if value == nil {
				continue
			}

			value = propertyConfig.ApplyVisibility(currentVisiblity, value)
			if value == nil {
				delete(extra.Fields, key)
			} else {
				extra.Fields[key] = value
			}
		}
	}

	return profile, nil
}

func (p *Providers) ValidateUserExtraData(pb *structpb.Struct) error {
	for key, value := range pb.Fields {
		propertyConfig, ok := p.Config.ExtraDataConfig[key]
		if !ok {
			return fmt.Errorf("%s: key not allowed", key)
		}

		if err := propertyConfig.Validate(value); err != nil {
			return fmt.Errorf("%s: %w", key, err)
		}
	}

	return nil
}
