package sms

import (
	"context"
	"fmt"

	twilio "github.com/kevinburke/twilio-go"
	"github.com/tierklinik-dobersberg/cis-idm/internal/config"
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

type NoopProvider struct{}

func NewNoopProvider() Sender {
	return new(NoopProvider)
}

func (*NoopProvider) Send(ctx context.Context, msg Message) error {
	return fmt.Errorf("sms provider is not configured")
}

// New creates a new SMSSender using acc.
func New(acc Account) (Sender, error) {
	client := twilio.NewClient(acc.AccountSid, acc.AccessToken, nil)

	return &sender{
		defaultFrom: acc.From,
		client:      client,
	}, nil
}

// SendTemplates renders a known template
func SendTemplate[T tmpl.Context](ctx context.Context, cfg config.Config, sender Sender, engine *tmpl.Engine, to []string, t tmpl.Known[T], args T) error {
	message, err := tmpl.RenderKnown(cfg, engine.SMS, t, args)
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
