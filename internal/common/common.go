package common

import (
	"strings"

	"github.com/tierklinik-dobersberg/cis-idm/internal/cache"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
)

type Service struct {
	repo  *repo.Repo
	cfg   config.Config
	cache cache.Cache
}

func New(repo *repo.Repo, cfg config.Config, cache cache.Cache) *Service {
	return &Service{
		repo:  repo,
		cfg:   cfg,
		cache: cache,
	}
}

func EnsureDisplayName(usr *models.User) {
	if usr.DisplayName != "" {
		return
	}

	if usr.FirstName != "" {
		usr.DisplayName = usr.FirstName
	}

	if usr.LastName != "" {
		usr.DisplayName = strings.TrimPrefix(usr.DisplayName+" "+usr.LastName, " ")
	}

	if usr.DisplayName == "" {
		usr.DisplayName = usr.Username
	}
}
