package providers

import (
	"context"

	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/email"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
	"github.com/mailersend/mailersend-go"
)

type mailersendClient struct {
	mail     *mailersend.Mailersend
	settings EmailSettings
	tracer   *tracing.Tracer
}

type EmailSettings struct {
	APIKey      string
	SenderName  string
	SenderEmail string
}

func NewMailerSend(settings EmailSettings, tracer *tracing.Provider) *mailersendClient {

	if settings.APIKey == "" {
		panic("Email API key is required")
	}

	return &mailersendClient{
		tracer:   tracer.NewTracer("provider-mailersend"),
		mail:     mailersend.NewMailersend(settings.APIKey),
		settings: settings,
	}
}

func (m *mailersendClient) Send(ctx context.Context, e ports.Email, template string) (err error) {
	var span *tracing.Span

	ctx, span = m.tracer.NewSpan(ctx, "send-email")
	defer func() {
		span.End(err)
	}()

	recipients := []mailersend.Recipient{
		{
			Email: e.To,
		},
	}

	from := mailersend.From{
		Name:  m.settings.SenderName,
		Email: m.settings.SenderEmail,
	}

	message := m.mail.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(e.Subject)
	message.SetHTML(template)

	_, err = m.mail.Email.Send(ctx, message)
	if err != nil {
		span.AddEvent("adapter-send-email-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	return nil
}
