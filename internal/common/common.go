package common

import (
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo"
)

type Service struct {
	repo *repo.Repo
	cfg  config.Config
}

func New(repo *repo.Repo, cfg config.Config) *Service {
	return &Service{
		repo: repo,
		cfg:  cfg,
	}
}
