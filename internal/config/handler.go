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
			"domain":       cfg.Server.Domain,
			"loginURL":     cfg.UserInterface.LoginRedirectURL,
			"siteName":     cfg.UserInterface.SiteName,
			"siteNameUrl":  cfg.UserInterface.SiteNameURL,
			"logoURL":      cfg.UserInterface.LogoURL,
			"registration": cfg.RegistrationMode,
			"features":     cfg.featureMap,
		}); err != nil {
			http.Error(w, "failed to encode config", http.StatusInternalServerError)

			return
		}

		if _, err := io.Copy(w, bytes.NewReader(buf.Bytes())); err != nil {
			logrus.WithError(err).Errorf("failed to send config response to client")
		}
	})
}
