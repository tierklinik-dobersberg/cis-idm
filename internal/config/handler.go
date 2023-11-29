package config

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

func NewConfigHandler(cfg Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)

		enc := json.NewEncoder(buf)
		enc.SetIndent("", "  ")

		if err := enc.Encode(map[string]any{
			"domain":                    cfg.Domain,
			"loginURL":                  cfg.LoginRedirectURL,
			"siteName":                  cfg.SiteName,
			"siteNameUrl":               cfg.SiteNameURL,
			"registrationRequiresToken": cfg.RegistrationRequiresToken,
			"features":                  cfg.featureMap,
			"logoURL":                   cfg.LogoURL,
		}); err != nil {
			http.Error(w, "failed to encode config", http.StatusInternalServerError)

			return
		}

		if _, err := io.Copy(w, bytes.NewReader(buf.Bytes())); err != nil {
			logrus.WithError(err).Errorf("failed to send config response to client")
		}
	})
}
