package config

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/apis/pkg/server"
	"golang.org/x/exp/slices"
)

var (
	ErrNotAllowed = errors.New("not allowed")
)

type Rule struct {
	IP        string `json:"ip"`
	Network   string `json:"network"`
	Token     string `json:"token"`
	Deny      bool   `json:"deny"`
	SubjectID string `json:"subject"`

	ip         net.IP
	network    *net.IPNet
	parseOnce  sync.Once
	parseError error
}

func (r *Rule) String() string {
	var parts []string
	if r.IP != "" {
		parts = append(parts, "ip:"+r.IP)
	}

	if r.Network != "" {
		parts = append(parts, "net:"+r.Network)
	}

	if r.Token != "" {
		parts = append(parts, "token:"+r.Token)
	}

	if r.Deny {
		parts = append(parts, "verdict:deny")
	} else {
		parts = append(parts, "verdict:allow")
	}

	parts = append(parts, "subject:"+r.SubjectID)

	return strings.Join(parts, " ")
}

type ForwardAuthEntry struct {
	Required *bool    `json:"required,omitempty" yaml:"required,omitempty"`
	URL      string   `json:"url" yaml:"url"`
	Methods  []string `json:"methods,omitempty" yaml:"methods,omitempty"`
	Rules    []*Rule  `json:"rules"`

	parsedUrlRegex *regexp.Regexp
	parseOnce      sync.Once
	parseError     error
}

// IsRequired returns true if authentication is required
// for this entry.
func (fae *ForwardAuthEntry) IsRequired() bool {
	if fae.Required == nil {
		return true
	}

	return *fae.Required
}

// Matches checks if fae matches url.
func (fae *ForwardAuthEntry) Matches(method, url string) (bool, error) {
	fae.parseOnce.Do(func() {
		fae.parsedUrlRegex, fae.parseError = regexp.CompilePOSIX(fae.URL)

		for idx := range fae.Methods {
			fae.Methods[idx] = strings.ToLower(fae.Methods[idx])
		}
	})

	if fae.parseError != nil {
		return false, fae.parseError
	}

	if len(fae.Methods) > 0 && !slices.Contains(fae.Methods, strings.ToLower(method)) {
		return false, nil
	}

	return fae.parsedUrlRegex.MatchString(url), nil
}

func (fae *ForwardAuthEntry) Allowed(req *http.Request) (string, bool, error) {
	merr := new(multierror.Error)

	subjectID := ""
	for _, r := range fae.Rules {
		log.L(req.Context()).Debugf("checking rule %s", r)

		matched, err := r.Matches(req)
		if err != nil {
			merr.Errors = append(merr.Errors, err)

			continue
		}

		if matched {
			// The rule explicitly denies access so there's no need to further
			// check other rules
			if r.Deny {
				return "", false, merr.ErrorOrNil()
			} else if subjectID == "" {
				subjectID = r.SubjectID

				// continue checking the other rules as one still
				// overwrite with Deny=true
			}
		}
	}

	allowed := false
	if !fae.IsRequired() {
		allowed = true
	} else {
		allowed = subjectID != ""
	}

	return subjectID, allowed, merr.ErrorOrNil()
}

func (r *Rule) Matches(req *http.Request) (bool, error) {
	r.parseOnce.Do(func() {
		if r.IP != "" {
			r.ip = net.ParseIP(r.IP)
			if r.ip == nil {
				r.parseError = fmt.Errorf("invalid IP address")
			}
		}

		if r.Network != "" {
			var err error

			_, r.network, err = net.ParseCIDR(r.Network)
			if err != nil {
				r.parseError = err
			}
		}
	})

	if r.parseError != nil {
		return false, r.parseError
	}

	l := log.L(req.Context())

	// get the real client IP address.
	clientIP := server.RealIPFromContext(req.Context())
	if clientIP == nil {
		l.Warnf("no real client ip associated with request context, usring RemoteAddr %q instead", req.RemoteAddr)

		sip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			return false, fmt.Errorf("failed to split client host:port (%q): %w", req.RemoteAddr, err)
		}

		clientIP = net.ParseIP(sip)
		if clientIP == nil {
			return false, fmt.Errorf("no client ip associated with request context and failed to parse IP %q", sip)
		}
	}

	// check if the IP matches
	if r.ip != nil {
		if clientIP.Equal(r.ip) {
			l.Infof("client IP matches rule.ip: %s", clientIP.String())
			return true, nil
		}

		l.Debugf("client IP %s does not match rule.ip %s", clientIP.String(), r.ip.String())
	}

	// check if the IP originates from a specified network.
	if r.network != nil {
		if r.network.Contains(clientIP) {
			l.Infof("client IP %s matches rule.network: %s", clientIP, r.network.String())

			return true, nil
		}

		l.Debugf("client IP %s does not match rule.network: %s", clientIP, r.network.String())
	}

	// check if there's a static token configured.
	if r.Token != "" {
		l.Debugf("checking for authorization bearer token")

		if h := req.Header.Get("Authorization"); h != "" {
			bearer, token, ok := strings.Cut(h, " ")

			if !ok {
				l.Debugf("invalid Authorization header %q", h)

				return false, nil
			}

			if strings.ToLower(bearer) != "bearer" {
				l.Debugf("invalid Authoriziation header %q", h)

				return false, nil
			}

			token = strings.TrimSpace(token)

			if r.Token == token {
				l.Info("client token matches static config token")

				return true, nil
			} else {
				l.Debugf("Authorization header set but token %q does not match", token)
			}
		}
	} else {
		l.Debugf("no authorization bearer token defined, request not allowed by rule")
	}

	return false, nil
}
