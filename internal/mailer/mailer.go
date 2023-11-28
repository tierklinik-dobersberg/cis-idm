package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"

	"github.com/ory/mail"
	"github.com/tierklinik-dobersberg/apis/pkg/log"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
)

type Account struct {
	// Host should holds the SMPT host name the mailer should
	// use to send mails.
	Host string
	// Port should holds the SMTP port the host is listening on.
	Port int
	// Username required for authentication.
	Username string
	// Password required for authentication.
	Password string
	// From defines the default sender to use for this account.
	From string
	// AllowInsecure can be set to true to disable TLS certificate
	// verification.
	AllowInsecure bool
	// UseSSL can be set to either true or false to force SSL to be enabled
	// or disabled. If not configured, SSL will be used for the default
	// SSL port.
	UseSSL *bool
}

type Message struct {
	From string
	To   []string
	Cc   []string
	Bcc  []string
}

type Mailer interface {
	DialAndSend(...*mail.Message) error
}

type mailer struct {
	*mail.Dialer

	defaultFrom string
}

type NoOpMailer struct{}

func (*NoOpMailer) DialAndSend(...*mail.Message) error {
	return fmt.Errorf("mail not configured")
}

func PrepareTemplate[T tmpl.Context](ctx context.Context, cfg config.Config, engine *tmpl.Engine, email Message, template tmpl.Known[T], args T) (*mail.Message, error) {
	var subjectTemplate tmpl.Known[T]
	subjectTemplate.Name = fmt.Sprintf("%s:subject", template.Name)

	subj, err := tmpl.RenderKnown(cfg, engine.Mail, subjectTemplate, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute subject template: %w", err)
	}

	msg := mail.NewMessage()
	msg.SetHeaders(map[string][]string{
		"From":    {email.From},
		"To":      email.To,
		"Cc":      email.Cc,
		"Bcc":     email.Bcc,
		"Subject": {subj},
	})

	log.L(ctx).Infof("Sending mail %q to %s, cc=%v and bcc=%v", subj, email.To, email.Cc, email.Bcc)

	msg.SetBodyWriter("text/html", func(w io.Writer) error {
		return tmpl.RenderKnownTo(cfg, engine.Mail, template, args, w)
	})

	return msg, nil
}

func SendTemplate[T tmpl.Context](ctx context.Context, cfg config.Config, engine *tmpl.Engine, m Mailer, email Message, template tmpl.Known[T], args T) error {
	// In dry-run mode, we replace the target address by fixed one
	if cfg.DryRun != nil && cfg.DryRun.MailTarget != "" {
		log.L(ctx).Infof("replacing e-mail receipients %v with %s in dry-run mode", email.To, cfg.DryRun.MailTarget)

		c := new(Message)
		*c = email

		c.To = []string{cfg.DryRun.MailTarget}
	}

	msg, err := PrepareTemplate(ctx, cfg, engine, email, template, args)
	if err != nil {
		return err
	}

	if err := m.DialAndSend(msg); err != nil {
		return fmt.Errorf("send-mail: %w", err)
	}

	return nil
}

// New returns a new mailer that sends mail through account.
func New(account Account) (*mailer, error) {
	dialer := mail.NewDialer(account.Host, account.Port, account.Username, account.Password)

	if account.AllowInsecure {
		if dialer.TLSConfig == nil {
			dialer.TLSConfig = new(tls.Config)
		}
		dialer.TLSConfig.InsecureSkipVerify = true
	}

	if account.UseSSL != nil {
		dialer.SSL = *account.UseSSL
	}

	return &mailer{
		Dialer:      dialer,
		defaultFrom: account.From,
	}, nil
}
