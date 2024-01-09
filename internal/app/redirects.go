package app

import (
	"context"
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/tierklinik-dobersberg/apis/pkg/log"
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

		for _, allowedDomain := range p.Config.Server.AllowedDomainRedirects {
			if strings.HasPrefix(allowedDomain, ".") {
				if strings.HasSuffix(u.Host, allowedDomain) {
					return u.String(), nil
				}
			}

			if u.Host == allowedDomain {
				return u.String(), nil
			}
		}

		log.L(ctx).Warnf("requested redirect to %s is not allowed (%v)", string(decoded), p.Config.Server.AllowedDomainRedirects)
	}

	return "", nil
}
