package helper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/axelspringer/swerve/log"
)

type logWriter struct {
	http.ResponseWriter
	status int
}

func (w *logWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *logWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	return n, err
}

// LoggingMiddleware logs request infos
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := logWriter{ResponseWriter: w}
		next.ServeHTTP(&sw, r)
		diff := time.Now().Sub(start)
		log.InfoWithFields(logrus.Fields{
			"method":  r.Method,
			"proto":   r.Proto,
			"code":    sw.status,
			"host":    r.Host,
			"path":    r.URL.Path,
			"qstring": r.URL.RawQuery,
			"took":    fmt.Sprintf("%.03fms", float64(diff.Microseconds())/1000),
			"ua":      r.UserAgent(),
			"remote":  r.RemoteAddr,
		}, "incoming request")
	})
}
