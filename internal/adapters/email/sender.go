package invite

import (
	"bytes"
	"context"
	libTemplate "text/template"

	tempEmail "github.com/Logstay/tech-challange-wherehouse-service/internal/adapters/email/html"
	"github.com/Logstay/tech-challange-wherehouse-service/internal/adapters/email/messages"
	"github.com/Logstay/tech-challange-wherehouse-service/internal/adapters/email/providers"
	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/email"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
)

type Email struct {
	client providers.Client
	tracer *tracing.Tracer
}

func NewEmail(client providers.Client, tracer *tracing.Provider) (sg *Email) {
	sg = &Email{
		tracer: tracer.NewTracer("email-Email"),
		client: client,
	}
	return
}

func (s *Email) SendEmail(ctx context.Context, e ports.Email) (err error) {
	var span *tracing.Span
	var internal ports.InternalParams
	internal.Email = e
	internal.Tmap = messages.GetMessages()

	ctx, span = s.tracer.NewSpan(ctx, "adapter-send-email")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("adapter-send-email-params", tracing.Attributes{
		"email": e,
	})

	searchedTemplate, err := tempEmail.GetTemplateHTML(e.SelectedTemplate)
	if err != nil {
		span.AddEvent("adapter-send-email-get-template-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	baseTemplate, err := s.templateBaseTransform(internal, string(searchedTemplate))
	if err != nil {
		span.AddEvent("adapter-send-email-template-base-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	bufTemplate, err := s.executeFInalTemplate(internal, baseTemplate)
	if err != nil {
		span.AddEvent("adapter-send-email-exec-template-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	err = s.client.Send(ctx, e, string(bufTemplate))
	if err != nil {
		span.AddEvent("adapter-send-email-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("adapter-send-email-sent", nil)
	return
}

func (s *Email) templateBaseTransform(data ports.InternalParams, template string) (resp string, err error) {
	buf := bytes.NewBuffer(nil)

	templateBase, err := tempEmail.GetTemplateHTML("base")
	if err != nil {
		return
	}

	h, err := libTemplate.New("").Parse(string(templateBase))
	if err != nil {
		return
	}

	t, err := h.Parse(template)
	if err != nil {
		return
	}

	err = t.Execute(buf, data)
	if err != nil {
		return
	}

	resp = buf.String()

	return
}

func (s *Email) executeFInalTemplate(data ports.InternalParams, finalTemplate string) (resp []byte, err error) {
	buf := bytes.NewBuffer(nil)

	t, err := libTemplate.New("").Parse(finalTemplate)
	if err != nil {
		return
	}

	err = t.Execute(buf, data)
	if err != nil {
		return
	}

	resp = buf.Bytes()

	return
}
