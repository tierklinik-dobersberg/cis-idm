package common

import (
	"context"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func (svc *Service) AddEmailAddressToUser(ctx context.Context, mailModel models.EMail) (*models.EMail, []models.EMail, error) {
	userID := mailModel.UserID

	if !svc.cfg.FeatureEnabled(config.FeatureEMails) {
		return nil, nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("email: %w", config.ErrFeatureDisabled))
	}

	addedMail, err := svc.repo.CreateUserEmail(ctx, mailModel)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to store new email address: %w", err)
	}

	mails, err := svc.repo.GetUserEmails(ctx, userID)
	if err != nil {
		return &addedMail, nil, fmt.Errorf("failed to get existing user emails: %w", err)
	}

	return &addedMail, mails, nil
}

func (svc *Service) DeleteEmailAddressFromUser(ctx context.Context, userID string, mailID string) ([]models.EMail, error) {
	if !svc.cfg.FeatureEnabled(config.FeatureEMails) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("email: %w", config.ErrFeatureDisabled))
	}

	log.L(ctx).WithField("email_id", mailID).Infof("deleting email address from user")

	if err := svc.repo.DeleteEMailFromUser(ctx, userID, mailID); err != nil {
		return nil, fmt.Errorf("failed to delete email from user: %w", err)
	}

	mails, err := svc.repo.GetUserEmails(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing user emails: %w", err)
	}

	return mails, nil
}

func (svc *Service) MarkEmailAsPrimary(ctx context.Context, userID string, mailID string) error {
	if !svc.cfg.FeatureEnabled(config.FeatureEMails) {
		return connect.NewError(connect.CodeUnavailable, fmt.Errorf("addresses: %w", config.ErrFeatureDisabled))
	}

	if err := svc.repo.MarkEmailAsPrimary(ctx, userID, mailID); err != nil {
		if errors.Is(err, stmts.ErrNoRowsAffected) {
			return connect.NewError(connect.CodeNotFound, fmt.Errorf("email address not found"))
		}

		return err
	}

	return nil
}
