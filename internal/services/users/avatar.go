package users

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/app"
	"github.com/vincent-petithory/dataurl"
)

func NewAvatarHandler(providers *app.Providers) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		userID := pathParts[len(pathParts)-1]

		user, err := providers.Datastore.GetUserByID(r.Context(), userID)
		if err != nil {
			http.Error(w, "failed to get user", http.StatusNotFound)
			return
		}

		if user.Avatar == "" {
			http.Error(w, "no avatar", http.StatusNotFound)
			return
		}

		if strings.HasPrefix(user.Avatar, "http") {
			http.Redirect(w, r, user.Avatar, http.StatusFound)
			return
		}

		du, err := dataurl.DecodeString(user.Avatar)
		if err != nil {
			http.Error(w, "invalid avatar data url", http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", du.ContentType())
		w.Header().Add("Cache-Control", "max-age=3600")
		if du.Encoding != "" {
			w.Header().Add("Encoding", du.Encoding)
		}

		w.WriteHeader(http.StatusOK)

		if _, err := io.Copy(w, bytes.NewReader(du.Data)); err != nil {
			log.L(r.Context()).With("error", err).Error("failed to send avatar data to connection")
		}
	})
}
