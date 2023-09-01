package logger

import (
	"context"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

const (
	TraceIDKey  = "X-Trace-ID"
	SpanIDKey   = "X-Span-ID"
	RequestKey  = "X-Request"
	ResponseKey = "X-Response"
	VersionKey  = "version"
	StackKey    = "stack"
)

var (
	version string
)

type Logger = logrus.Logger

type Entry = logrus.Entry

// SetLevel logrus iota value
func SetLevel(level int) {
	logrus.SetLevel(logrus.Level(level))
}

// SetFormatter
func SetFormatter(format string) {
	switch format {
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	default:
		logrus.SetFormatter(new(logrus.TextFormatter))
	}
}

// SetReportCaller
func SetReportCaller(include bool) {
	logrus.SetReportCaller(include)
}

// SetOutput
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

// SetVersion
func SetVersion(v string) {
	version = v
}

func GetWrite() *io.PipeWriter {
	return logrus.StandardLogger().Writer()
}

type (
	traceIDKey  struct{}
	spanIDKey   struct{}
	requestKey  struct{}
	responseKey struct{}
	stackKey    struct{}
)

// NewTraceIDContext
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// FromTraceIDContext
func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewSpanIDContext
func NewSpanIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, spanIDKey{}, traceID)
}

// FromSpanIDContext
func FromSpanIDContext(ctx context.Context) string {
	v := ctx.Value(spanIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewRequestContext
func NewRequestContext(ctx context.Context, request string) context.Context {
	return context.WithValue(ctx, requestKey{}, request)
}

// FromRequestContext
func FromRequestContext(ctx context.Context) string {
	v := ctx.Value(requestKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewResponseContext
func NewResponseContext(ctx context.Context, response string) context.Context {
	return context.WithValue(ctx, responseKey{}, response)
}

// FromResponseContext
func FromResponseContext(ctx context.Context) string {
	v := ctx.Value(responseKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewStackContext
func NewStackContext(ctx context.Context, stack error) context.Context {
	return context.WithValue(ctx, stackKey{}, stack)
}

// FromStackContext
func FromStackContext(ctx context.Context) error {
	v := ctx.Value(stackKey{})
	if v != nil {
		if s, ok := v.(error); ok {
			return s
		}
	}
	return nil
}

// WithContext Use context create entry
func WithContext(ctx context.Context) *Entry {
	if ctx == nil {
		ctx = context.Background()
	}

	fields := map[string]interface{}{
		VersionKey: version,
	}

	if v := FromTraceIDContext(ctx); v != "" {
		fields[TraceIDKey] = v
	}

	if v := FromSpanIDContext(ctx); v != "" {
		fields[SpanIDKey] = v
	}

	if v := FromStackContext(ctx); v != nil {
		fields[StackKey] = fmt.Sprintf("%+v", v)
	}

	if v := FromRequestContext(ctx); v != "" {
		fields[RequestKey] = v
	}

	if v := FromResponseContext(ctx); v != "" {
		fields[ResponseKey] = v
	}

	return logrus.WithContext(ctx).WithFields(fields)
}

// Define logrus alias
var (
	Tracef = logrus.Tracef
	Debugf = logrus.Debugf
	Infof  = logrus.Infof
	Warnf  = logrus.Warnf
	Errorf = logrus.Errorf
	Fatalf = logrus.Fatalf
	Panicf = logrus.Panicf
	Printf = logrus.Printf
)
