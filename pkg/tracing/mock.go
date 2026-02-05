package tracing

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

type TraceMock struct {
	embedded.Tracer
	mock.Mock
}

func (t *TraceMock) Start(parentCtx context.Context,
	name string, ops ...trace.SpanStartOption) (context.Context, trace.Span) {
	args := t.Called(parentCtx, name, ops)
	return args.Get(0).(context.Context), args.Get(1).(*TraceMockSpan)
}

type TraceMockSpan struct {
	embedded.Span
	mock.Mock
}

func (t *TraceMockSpan) End(options ...trace.SpanEndOption) {
}

func (t *TraceMockSpan) AddEvent(name string, options ...trace.EventOption) {

}

func (t *TraceMockSpan) AddLink(l trace.Link) {}

func (t *TraceMockSpan) IsRecording() bool {
	return false
}

func (t *TraceMockSpan) RecordError(err error, options ...trace.EventOption) {}

func (t *TraceMockSpan) SpanContext() trace.SpanContext {
	return trace.SpanContext{}
}

func (t *TraceMockSpan) SetStatus(code codes.Code, description string) {}

func (t *TraceMockSpan) SetName(name string) {}

func (t *TraceMockSpan) SetAttributes(kv ...attribute.KeyValue) {}

func (t *TraceMockSpan) TracerProvider() trace.TracerProvider {
	return &TraceProviderMock{}
}

type TraceProviderMock struct {
	embedded.TracerProvider
}

func (t *TraceProviderMock) Tracer(name string, options ...trace.TracerOption) trace.Tracer {
	return &TraceMock{}
}
