package browser

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger(t *testing.T) *zap.Logger {
	sync := zapcore.AddSync(tlogWriter(t.Log))
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		sync,
		zap.DebugLevel,
	)
	logger := zap.New(core)
	t.Cleanup(func() {
		logger.Sync()
	})
	return logger
}

type tlogWriter func(...interface{})

func (w tlogWriter) Write(p []byte) (int, error) {
	w(string(p))
	return len(p), nil
}
