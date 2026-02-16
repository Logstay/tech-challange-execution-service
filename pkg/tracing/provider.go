package tracing

import (
	"context"

	"github.com/Logstay/tech-challange-execution-service/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Provider struct {
	*sdktrace.TracerProvider
}

func newResource(attrs ...attribute.KeyValue) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		attrs...,
	)
}

func NewTracingProvider(cfg *config.Config) (p *Provider, err error) {
	res := newResource(
		semconv.ServiceNameKey.String(cfg.Tracer.ServiceName),
		semconv.ServiceVersionKey.String(cfg.Tracer.ServiceVersion),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	p = &Provider{tracerProvider}

	return
}

func (p Provider) NewTracer(name string) *Tracer {
	otelTracer := p.Tracer(name)
	return &Tracer{otelTracer}
}

func (p Provider) Close(ctx context.Context) {
	err := p.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}
