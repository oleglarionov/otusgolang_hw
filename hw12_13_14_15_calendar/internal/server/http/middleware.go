package internalhttp

import (
	"fmt"
	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/common"
	"net"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(responseWriter http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: responseWriter}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler, l common.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := newLoggingResponseWriter(w)

		start := time.Now()
		next.ServeHTTP(lrw, r)
		elapsed := time.Since(start)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			l.Error("error parsing ip: " + err.Error())
			return
		}

		l.Info(fmt.Sprintf("%s [%s] %s %s %d %s %s",
			ip,
			start.Format("02/Jan/2006:03:04:05 Z0700"),
			r.Method,
			r.RequestURI,
			lrw.statusCode,
			elapsed,
			r.UserAgent(),
		))
	})
}
