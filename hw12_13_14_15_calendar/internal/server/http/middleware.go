package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/oleglarionov/otusgolang_hw/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware(next http.Handler, l logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
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
			200,
			elapsed,
			r.UserAgent(),
		))
	})
}
