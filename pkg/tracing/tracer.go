package tracing

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	trace.Tracer
}

func (t Tracer) NewSpan(parentCtx context.Context, name string) (ctx context.Context, span *Span) {
	ctx, otelSpan := t.Start(parentCtx, name)
	span = &Span{Span: otelSpan}

	if txn := newrelic.FromContext(parentCtx); txn != nil {
		segment := txn.StartSegment(name)
		span.nrSegment = segment
	}

	return
}

func (t Tracer) GetSpanFromContext(ctx context.Context) *Span {
	return &Span{Span: trace.SpanFromContext(ctx)}
}
