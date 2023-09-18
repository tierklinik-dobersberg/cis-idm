package tmpl

import (
	"bytes"
	"fmt"
	htmlTemplate "html/template"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	textTemplate "text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/ory/mail"
	"github.com/sirupsen/logrus"
	idmv1 "github.com/tierklinik-dobersberg/apis/gen/go/tkd/idm/v1"
	"github.com/tierklinik-dobersberg/apis/pkg/overlayfs"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/repo/models"
	"github.com/vincent-petithory/dataurl"
	"golang.org/x/exp/slices"
)

type RenderContext struct {
	values map[string]any
}

func (rc *RenderContext) Set(key string, value any) {
	rc.values[key] = value
}

func (rc *RenderContext) Get(key string) any {
	return rc.values[key]
}

func NewRenderContext() *RenderContext {
	return &RenderContext{
		values: make(map[string]any),
	}
}

type Engine struct {
	SMS  TemplateEngine
	Mail TemplateEngine
}

func New(fileSystems ...fs.FS) (*Engine, error) {
	mergedFS := overlayfs.NewFS(append(fileSystems, builtin)...)

	sms, err := NewTextEngine(mergedFS, KindSMS)
	if err != nil {
		return nil, fmt.Errorf("failed to create sms template engine: %w", err)
	}
	mail, err := NewHTMLEngine(mergedFS, KindMail)
	if err != nil {
		return nil, fmt.Errorf("failed to create mail template engine: %w", err)
	}

	return &Engine{
		SMS:  sms,
		Mail: mail,
	}, nil
}

type TemplateEngine interface {
	ExecuteTemplate(wr io.Writer, name string, data any) error
}

func NewTextEngine(fs fs.FS, kind Kind) (TemplateEngine, error) {
	t := textTemplate.New("")
	fm := textTemplate.FuncMap(PrepareFunctionMap())

	t.Funcs(fm)

	t, err := t.ParseFS(fs, filepath.Join("templates", string(kind), "*.tmpl"))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewHTMLEngine(fs fs.FS, kind Kind) (TemplateEngine, error) {
	t := htmlTemplate.New("")
	fm := htmlTemplate.FuncMap(PrepareFunctionMap())

	t.Funcs(fm)

	t, err := t.ParseFS(fs, filepath.Join("templates", string(kind), "*.html"))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func RenderKnown[T Context](cfg config.Config, engine TemplateEngine, known Known[T], args T) (string, error) {
	var buf = new(strings.Builder)

	if err := RenderKnownTo(cfg, engine, known, args, buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func RenderKnownTo[T Context](cfg config.Config, engine TemplateEngine, known Known[T], args T, target io.Writer) error {
	args.SetPublicURL(cfg.PublicURL)
	args.SetSiteName(cfg.SiteName)
	args.SetSiteURL(cfg.SiteNameURL)

	if err := engine.ExecuteTemplate(target, known.Name, args); err != nil {
		return err
	}

	return nil
}

// AddToMap - add src's entries to dst
func AddToMap(dst, src map[string]any) {
	for k, v := range src {
		dst[k] = v
	}
}

func PrepareFunctionMap() map[string]any {
	m := make(map[string]any)

	AddToMap(m, sprig.GenericFuncMap())
	AddToMap(m, customMap)

	return m
}

var customMap = map[string]any{
	"primaryEmail": func(input *idmv1.Profile) string {
		if pm := input.User.PrimaryMail; pm != nil {
			return pm.Address
		}

		return ""
	},
	"primaryPhone": func(input *idmv1.Profile) string {
		if pp := input.User.PrimaryPhoneNumber; pp != nil {
			return pp.Number
		}

		return ""
	},
	"displayName": func(input any) string {
		if profile, ok := input.(*idmv1.Profile); ok {
			if profile.User.DisplayName != "" {
				return profile.User.DisplayName
			}

			return profile.User.Username
		}

		if model, ok := input.(models.User); ok {
			if model.DisplayName != "" {
				return model.DisplayName
			}

			return model.Username
		}

		panic("expected *idmv1.Profile or models.User")
	},
	"userAvatar": func(input any, ctx *RenderContext) htmlTemplate.URL {
		m, ok := ctx.Get("mail").(*mail.Message)
		if !ok {
			return ""
		}

		var (
			userID    string
			avatarURL string
		)

		if profile, ok := input.(*idmv1.Profile); ok {
			userID = profile.User.Id
			avatarURL = profile.User.Avatar
		} else if model, ok := input.(models.User); ok {
			userID = model.ID
			avatarURL = model.Avatar
		}

		attachments, ok := ctx.Get("attachments").([]string)
		if !ok || !slices.Contains(attachments, "cid:"+userID) {
			attachUserAvatar(m, avatarURL, userID)
		} else {
			attachments = append(attachments, "cid:"+userID)
			ctx.Set("attachments", attachments)
		}

		return htmlTemplate.URL("cid:" + userID)
	},
}

func attachUserAvatar(m *mail.Message, avatarURL string, userID string) {
	ext, avatar, ct, err := getUserAvatar(avatarURL)
	if err != nil {
		logrus.Errorf("failed to get sender avatar as dataurl: %s", err)

		return
	}

	m.EmbedReader(userID+ext, bytes.NewReader(avatar), mail.SetHeader(map[string][]string{
		"Content-Type": {ct + "; name=" + fmt.Sprintf("%q", userID+ext)},
		"Content-ID":   {"<" + userID + ">"},
	}))
}

func getUserAvatar(avatar string) (string, []byte, string, error) {
	var du *dataurl.DataURL

	if strings.HasPrefix(avatar, "http") {
		res, err := http.Get(avatar)
		if err != nil {
			return "", nil, "", fmt.Errorf("failed to load avatar: %w", err)
		}
		defer res.Body.Close()

		content, err := io.ReadAll(res.Body)
		if err != nil {
			return "", nil, "", fmt.Errorf("failed to read body: %w", err)
		}

		du = dataurl.New(content, res.Header.Get("Content-Type"))
	} else {
		var err error
		du, err = dataurl.DecodeString(avatar)
		if err != nil {
			return "", nil, "", fmt.Errorf("failed to decode dataurl: %w", err)
		}
	}

	exts, err := mime.ExtensionsByType(du.MediaType.ContentType())
	if err != nil || len(exts) == 0 {
		exts = []string{".img"}
	}

	return exts[0], du.Data, du.MediaType.ContentType(), nil
}
