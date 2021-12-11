package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger) func(http.Handler) http.Handler {
	return middleware.RequestLogger(&zapRequestLogger{
		logger: logger,
	})
}

type zapRequestLogger struct {
	logger *zap.Logger
}

var _ middleware.LogFormatter = (*zapRequestLogger)(nil)

func (l *zapRequestLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	logger := l.logger.With(
		zap.String("method", r.Method),
		zap.String("uri", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)),
	)

	return &zapLogEntry{
		logger: logger,
	}
}

type zapLogEntry struct {
	logger *zap.Logger
}

var _ middleware.LogEntry = (*zapLogEntry)(nil)

func (l *zapLogEntry) Write(status int, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.logger.Info("request",
		zap.Int("status", status),
		zap.Duration("latency", elapsed),
	)
}

func (l *zapLogEntry) Panic(v interface{}, stack []byte) {
	l.logger.Panic("request",
		zap.ByteString("stack", stack),
		zap.String("panic", fmt.Sprintf("%+v", v)),
	)
}
