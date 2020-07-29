package helper

import (
	"net/http"
	"time"

	"github.com/axelspringer/swerve/log"
)

type logWriter struct {
	http.ResponseWriter
	status int
	length int
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
	w.length += n
	return n, err
}

// LoggingMiddleware logs request infos
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := logWriter{ResponseWriter: w}
		next.ServeHTTP(&sw, r)
		diff := time.Now().Sub(start)
		milli := int64(diff / time.Millisecond)
		micro := int64(diff / time.Microsecond)
		ts := time.Now().Format("02/Jan/2006 03:04:05")
		log.Infof(`ts="%s" method="%s" proto="%s" code="%d" host="%s" path="%s" qstring="%v" took="%d.%03dms" ua="%s"`, ts, r.Method, r.Proto, sw.status, r.Host, r.URL.Path, r.URL.RawQuery, milli, micro, r.UserAgent())
	})
}
