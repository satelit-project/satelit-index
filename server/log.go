package server

import (
	"bytes"
	"net/http"
	"satelit-project/satelit-index/logging"

	"github.com/gorilla/handlers"
	"go.uber.org/zap"
)

type loggerWriter struct {
	logFunc func(args ...interface{})
}

func NewLoggingHandler(logger *zap.SugaredLogger, h http.Handler) http.Handler {
	if logger == nil {
		logger = logging.DefaultLogger()
	}

	writer := loggerWriter{logger.Info}
	return handlers.LoggingHandler(&writer, h)
}

func (l *loggerWriter) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	l.logFunc(string(p))
	return len(p), nil
}
