package config

import (
	"regexp"
	"strings"
	"sync"

	"golang.org/x/exp/slices"
)

type StaticToken struct {
	Tokens    string   `json:"token"`
	SubjectID string   `json:"subject"`
	Roles     []string `json:"roles"`
}

type ForwardAuthEntry struct {
	Required *bool         `json:"required,omitempty" yaml:"required,omitempty"`
	URL      string        `json:"url" yaml:"url"`
	Methods  []string      `json:"methods,omitempty" yaml:"methods,omitempty"`
	Tokens   []StaticToken `json:"staticTokens"`

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
