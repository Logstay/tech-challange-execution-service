package tracing

import (
	"encoding/json"
	"fmt"

	"github.com/newrelic/go-agent/v3/newrelic"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Attributes map[string]any

type Span struct {
	trace.Span
	nrSegment *newrelic.Segment
}

func (s Span) AddEvent(name string, extra ...Attributes) {
	var attributes []attribute.KeyValue

	for _, m := range extra {
		for k, v := range m {
			s.processAttribute(k, v, &attributes)
		}
	}

	s.Span.AddEvent(name, trace.WithAttributes(attributes...))
}

func (s Span) End(errs ...error) {
	for _, err := range errs {
		if err != nil {
			s.RecordError(err)
		}
	}

	if s.nrSegment != nil {
		s.nrSegment.End()
	}

	s.Span.End()
}

func (s Span) processAttribute(k string, v any, attributes *[]attribute.KeyValue) {
	switch k {
	case "eventType", "http.status_code":
		if strValue, err := s.toString(v); err == nil {
			s.SetAttributes(attribute.String(k, strValue))
			if s.nrSegment != nil {
				s.nrSegment.AddAttribute(k, strValue)
			}
		}
	case "err":
		s.SetAttributes(attribute.String("error", "true"))
		if strValue, err := s.toString(v); err == nil {
			*attributes = append(*attributes, attribute.String(k, strValue))
			if s.nrSegment != nil {
				s.nrSegment.AddAttribute("error.message", strValue)
			}
		}
	default:
		jsonValue, err := s.ensureJSON(k, v)
		if err != nil {
			return
		}
		cleanedValue, err := removeSensitiveFields(jsonValue, "password")
		if err != nil {
			return
		}
		*attributes = append(*attributes, attribute.String(k, cleanedValue))
		if s.nrSegment != nil {
			s.nrSegment.AddAttribute(k, cleanedValue)
		}
	}
}

func (s Span) ensureJSON(k string, v any) (string, error) {
	vToMarshal := v

	if str, ok := v.(string); ok && !json.Valid([]byte(str)) {
		vToMarshal = map[string]any{k: str}
	} else if str, ok := v.(string); ok {
		return str, nil
	}

	bytes, err := json.Marshal(vToMarshal)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (s Span) toString(v any) (string, error) {
	switch value := v.(type) {
	case string:
		return value, nil
	case int, int32, int64, float32, float64, bool:
		return fmt.Sprintf("%v", value), nil
	case error:
		return value.Error(), nil
	default:
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(jsonValue), nil
	}
}

func removeSensitiveFields(jsonStr string, fields ...string) (string, error) {
	var data map[string]any

	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", err
	}

	for _, field := range fields {
		delete(data, field)
	}

	cleanedBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(cleanedBytes), nil
}
