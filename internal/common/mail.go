package common

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"github.com/gofrs/uuid"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

func (svc *Service) AddEmailAddressToUser(ctx context.Context, mailModel repo.UserEmail) (*repo.UserEmail, []repo.UserEmail, error) {
	userID := mailModel.UserID

	if !svc.cfg.FeatureEnabled(config.FeatureEMails) {
		return nil, nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("email: %w", config.ErrFeatureDisabled))
	}

	if mailModel.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, nil, err
		}

		mailModel.ID = id.String()
	}

	addedMail, err := svc.repo.CreateEMail(ctx, repo.CreateEMailParams{
		ID:       mailModel.ID,
		UserID:   mailModel.UserID,
		Address:  mailModel.Address,
		Verified: mailModel.Verified,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to store new email address: %w", err)
	}

	mails, err := svc.repo.GetEmailsForUserByID(ctx, userID)
	if err != nil {
		return &addedMail, nil, fmt.Errorf("failed to get existing user emails: %w", err)
	}

	return &addedMail, mails, nil
}

func (svc *Service) DeleteEmailAddressFromUser(ctx context.Context, userID string, mailID string) ([]repo.UserEmail, error) {
	if !svc.cfg.FeatureEnabled(config.FeatureEMails) {
		return nil, connect.NewError(connect.CodeUnavailable, fmt.Errorf("email: %w", config.ErrFeatureDisabled))
	}

	log.L(ctx).WithField("email_id", mailID).Infof("deleting email address from user")

	if rows, err := svc.repo.DeleteEMailFromUser(ctx, repo.DeleteEMailFromUserParams{
		ID:     mailID,
		UserID: userID,
	}); err == nil {
		if rows == 0 {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("mail with id %q not found", mailID))
		}
	} else {
		return nil, fmt.Errorf("failed to delete email from user: %w", err)
	}

	mails, err := svc.repo.GetEmailsForUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing user emails: %w", err)
	}

	return mails, nil
}

func (svc *Service) MarkEmailAsPrimary(ctx context.Context, userID string, mailID string) error {
	if !svc.cfg.FeatureEnabled(config.FeatureEMails) {
		return connect.NewError(connect.CodeUnavailable, fmt.Errorf("addresses: %w", config.ErrFeatureDisabled))
	}

	if rows, err := svc.repo.MarkEmailAsPrimary(ctx, repo.MarkEmailAsPrimaryParams{
		ID:     mailID,
		UserID: userID,
	}); err == nil {
		if rows == 0 {
			return connect.NewError(connect.CodeNotFound, fmt.Errorf("mail with id %q not found", mailID))
		}
	} else {
		return err
	}

	return nil
}
