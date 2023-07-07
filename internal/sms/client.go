package sms

import (
	"context"
	"fmt"

	twilio "github.com/kevinburke/twilio-go"
	"github.com/tierklinik-dobersberg/cis-idm/internal/tmpl"
)

type Account struct {
	From        string
	AccountSid  string
	AccessToken string
}

type Message struct {
	From string
	To   []string
	Body string
}

// Sender can send messages using the Twilio Programmable
// SMS interface.
type Sender interface {
	// Send sends msg and returns any error encountered.
	Send(ctx context.Context, msg Message) error
}

// New creates a new SMSSender using acc.
func New(acc Account, engine tmpl.TemplateEngine) (Sender, error) {
	client := twilio.NewClient(acc.AccountSid, acc.AccessToken, nil)

	return &sender{
		defaultFrom: acc.From,
		client:      client,
		engine:      engine,
	}, nil
}

// SendTemplates renders a known template
func SendTemplate[T tmpl.Context](ctx context.Context, sender Sender, engine tmpl.TemplateEngine, to []string, t tmpl.Known[T], args T) error {
	message, err := tmpl.RenderKnown(engine, t, args)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	msg := Message{
		To:   to,
		Body: message,
	}

	return sender.Send(ctx, msg)
}

type sender struct {
	client      *twilio.Client
	defaultFrom string
	engine      tmpl.TemplateEngine
}

func (s *sender) Send(ctx context.Context, msg Message) error {
	if msg.From == "" {
		msg.From = s.defaultFrom
	}

	for _, to := range msg.To {
		_, err := s.client.Messages.SendMessage(msg.From, to, msg.Body, nil)
		if err != nil {
			return fmt.Errorf("sending message to %s: %w", to, err)
		}
	}

	return nil
}
