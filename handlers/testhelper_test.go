package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type tlogWriter func(...interface{})

func (w tlogWriter) Write(p []byte) (int, error) {
	w(string(p))
	return len(p), nil
}

func testLogger(t *testing.T) *zap.Logger {
	sync := zapcore.AddSync(tlogWriter(t.Log))
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		sync,
		zap.DebugLevel,
	)
	return zap.New(core)
}

func get(handler http.Handler, path string) (int, http.Header, string) {
	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)
	result := response.Result()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}

	return result.StatusCode, result.Header, strings.TrimSpace(string(body))
}
