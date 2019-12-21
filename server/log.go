package server

import (
	"net/http"
	"time"

	"shitty.moe/satelit-project/satelit-index/logging"
)

// ResponseWriter which keeps track of response status code.
type loggingWriter struct {
	inner      http.ResponseWriter
	statusCode int
}

func (w *loggingWriter) Header() http.Header {
	return w.inner.Header()
}

func (w *loggingWriter) Write(c []byte) (int, error) {
	return w.inner.Write(c)
}

func (w *loggingWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.inner.WriteHeader(statusCode)
}

// Returns Handler which logs every request with provided logger.
func LogRequest(h http.Handler, log *logging.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		lw := &loggingWriter{inner: w}
		h.ServeHTTP(lw, r)

		if lw.statusCode == 0 {
			lw.statusCode = 200
		}

		log.Infof("%d\t%s\t%s\t%s", lw.statusCode, r.Method, r.RequestURI, time.Since(now))
	})
}
