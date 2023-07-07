package app

import (
	"context"
	"encoding/base64"
	"net/url"

	"github.com/tierklinik-dobersberg/cis-idm/internal/middleware"
	"golang.org/x/exp/slices"
)

func (p *Providers) HandleRequestedRedirect(ctx context.Context, requestedRedirect string) (string, error) {
	if requestedRedirect != "" {
		decoded, err := base64.URLEncoding.DecodeString(requestedRedirect)
		if err != nil {
			return "", err
		}

		u, err := url.Parse(string(decoded))
		if err != nil {
			return "", err
		}

		if slices.Contains(p.Config.AllowedDomainRedirects, u.Host) {
			middleware.L(ctx).Infof("redirecting user to %s", u.String())
			return u.String(), nil

		} else {
			middleware.L(ctx).Warnf("requested redirect to %s is not allowed", string(decoded))
		}
	}

	return "", nil
}
