package tmpl

import (
	"fmt"
	htmlTemplate "html/template"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	textTemplate "text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/tierklinik-dobersberg/cis-idm/internal/overlayfs"
)

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
	t, err := textTemplate.New("").ParseFS(fs, filepath.Join(string(kind), "*.tmpl"))
	if err != nil {
		return nil, err
	}

	fm := textTemplate.FuncMap{}

	addToMap(fm, sprig.GenericFuncMap())

	t.Funcs(fm)

	return t, nil
}

func NewHTMLEngine(fs fs.FS, kind Kind) (TemplateEngine, error) {
	t, err := htmlTemplate.New("").ParseFS(fs, filepath.Join(string(kind), "*.tmpl"))
	if err != nil {
		return nil, err
	}

	fm := htmlTemplate.FuncMap{}

	addToMap(fm, sprig.GenericFuncMap())

	t.Funcs(fm)

	return t, nil
}

func RenderKnown[T Context](engine TemplateEngine, known Known[T], args T) (string, error) {
	var buf = new(strings.Builder)

	if err := engine.ExecuteTemplate(buf, known.Name, args); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// addToMap - add src's entries to dst
func addToMap(dst, src map[string]interface{}) {
	for k, v := range src {
		dst[k] = v
	}
}
