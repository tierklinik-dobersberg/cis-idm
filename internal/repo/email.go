package repo

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/stmts"
)

func (repo *Repo) GetUserEmails(ctx context.Context, userID string) ([]models.EMail, error) {
	return Query(ctx, stmts.GetEmailsForUserByID, repo.Conn, map[string]any{
		"user_id": userID,
	})
}

func (repo *Repo) GetUserPrimaryMail(ctx context.Context, userID string) (models.EMail, error) {
	return QueryOne(ctx, stmts.GetPrimaryEmailForUserByID, repo.Conn, map[string]any{
		"user_id": userID,
	})
}

func (repo *Repo) CreateUserEmail(ctx context.Context, mail models.EMail) (models.EMail, error) {
	if mail.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return mail, err
		}

		mail.ID = id.String()
	}

	if err := stmts.CreateEMail.Write(ctx, repo.Conn, mail); err != nil {
		return mail, err
	}

	return mail, nil
}

func (repo *Repo) DeleteEMailFromUser(ctx context.Context, userID string, mailID string) error {
	return stmts.DeleteEMailFromUser.Write(ctx, repo.Conn, map[string]any{
		"user_id": userID,
		"id":      mailID,
	})
}

func (repo *Repo) MarkEmailAsPrimary(ctx context.Context, userID string, mailID string) error {
	return stmts.MarkEmailAsPrimary.Write(ctx, repo.Conn, map[string]any{
		"user_id": userID,
		"id":      mailID,
	})
}
