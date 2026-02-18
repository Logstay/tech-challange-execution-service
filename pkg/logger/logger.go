package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
	LevelFatal LogLevel = "fatal"
)

type StructuredLog struct {
	Timestamp  string         `json:"timestamp"`
	Level      LogLevel       `json:"level"`
	Message    string         `json:"message"`
	RequestID  string         `json:"request_id,omitempty"`
	TraceID    string         `json:"trace_id,omitempty"`
	UserID     string         `json:"user_id,omitempty"`
	Method     string         `json:"method,omitempty"`
	Path       string         `json:"path,omitempty"`
	StatusCode int            `json:"status_code,omitempty"`
	Latency    string         `json:"latency,omitempty"`
	Error      string         `json:"error,omitempty"`
	Service    string         `json:"service"`
	Extra      map[string]any `json:"extra,omitempty"`
}

type Logger interface {
	Debug(msg string, extra ...map[string]any)
	Info(msg string, extra ...map[string]any)
	Warn(msg string, extra ...map[string]any)
	Error(msg string, err error, extra ...map[string]any)
	Fatal(msg string, err error, extra ...map[string]any)
	WithContext(ctx context.Context) Logger
	WithRequestID(requestID string) Logger
	WithTraceID(traceID string) Logger
	WithUserID(userID string) Logger
	WithError(err error) Logger
}

type logger struct {
	service   string
	requestID string
	traceID   string
	userID    string
	logrus    *logrus.Logger
}

func NewLogger(service string, logrusInstance *logrus.Logger) Logger {
	return &logger{
		service: service,
		logrus:  logrusInstance,
	}
}

func (l *logger) Debug(msg string, extra ...map[string]any) {
	l.log(LevelDebug, msg, nil, extra...)
}

func (l *logger) Info(msg string, extra ...map[string]any) {
	l.log(LevelInfo, msg, nil, extra...)
}

func (l *logger) Warn(msg string, extra ...map[string]any) {
	l.log(LevelWarn, msg, nil, extra...)
}

func (l *logger) Error(msg string, err error, extra ...map[string]any) {
	l.log(LevelError, msg, err, extra...)
}

func (l *logger) Fatal(msg string, err error, extra ...map[string]any) {
	l.log(LevelFatal, msg, err, extra...)
	os.Exit(1)
}

func (l *logger) WithContext(ctx context.Context) Logger {
	newLogger := &logger{
		service:   l.service,
		requestID: l.requestID,
		traceID:   l.traceID,
		userID:    l.userID,
		logrus:    l.logrus,
	}

	if ginCtx, ok := ctx.(*gin.Context); ok {
		if reqID, exists := ginCtx.Get("request_id"); exists {
			newLogger.requestID = reqID.(string)
		}
		if traceID, exists := ginCtx.Get("trace_id"); exists {
			newLogger.traceID = traceID.(string)
		}
		if userID, exists := ginCtx.Get("user_id"); exists {
			newLogger.userID = fmt.Sprintf("%v", userID)
		}
	}

	return newLogger
}

func (l *logger) WithRequestID(requestID string) Logger {
	return &logger{
		service:   l.service,
		requestID: requestID,
		traceID:   l.traceID,
		userID:    l.userID,
		logrus:    l.logrus,
	}
}

func (l *logger) WithTraceID(traceID string) Logger {
	return &logger{
		service:   l.service,
		requestID: l.requestID,
		traceID:   traceID,
		userID:    l.userID,
		logrus:    l.logrus,
	}
}

func (l *logger) WithUserID(userID string) Logger {
	return &logger{
		service:   l.service,
		requestID: l.requestID,
		traceID:   l.traceID,
		userID:    userID,
		logrus:    l.logrus,
	}
}

func (l *logger) WithError(err error) Logger {
	return &logger{
		service:   l.service,
		requestID: l.requestID,
		traceID:   l.traceID,
		userID:    l.userID,
		logrus:    l.logrus,
	}
}

func (l *logger) log(level LogLevel, msg string, err error, extra ...map[string]any) {
	fields := logrus.Fields{
		"service": l.service,
	}

	if l.requestID != "" {
		fields["request_id"] = l.requestID
	}
	if l.traceID != "" {
		fields["trace_id"] = l.traceID
	}
	if l.userID != "" {
		fields["user_id"] = l.userID
	}

	if len(extra) > 0 && extra[0] != nil {
		for k, v := range extra[0] {
			fields[k] = v
		}
	}

	entry := l.logrus.WithFields(fields)

	if err != nil {
		entry = entry.WithError(err)
	}

	switch level {
	case LevelDebug:
		entry.Debug(msg)
	case LevelInfo:
		entry.Info(msg)
	case LevelWarn:
		entry.Warn(msg)
	case LevelError:
		entry.Error(msg)
	case LevelFatal:
		entry.Fatal(msg)
	}
}

func RequestLogger(service string, logrusInstance *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		startTime := time.Now()

		c.Next()

		latency := time.Since(startTime)

		fields := logrus.Fields{
			"service":     service,
			"request_id":  requestID,
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status_code": c.Writer.Status(),
			"latency":     latency.String(),
		}

		if traceID, exists := c.Get("trace_id"); exists {
			fields["trace_id"] = traceID.(string)
		}

		if userID, exists := c.Get("user_id"); exists {
			fields["user_id"] = fmt.Sprintf("%v", userID)
		}

		entry := logrusInstance.WithFields(fields)

		if c.Writer.Status() >= 400 {
			if len(c.Errors) > 0 {
				entry.WithError(c.Errors.Last()).Error("HTTP Request")
			} else {
				entry.Error("HTTP Request")
			}
		} else {
			entry.Info("HTTP Request")
		}
	}
}
