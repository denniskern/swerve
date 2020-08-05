package helper

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/axelspringer/swerve/database"

	"github.com/axelspringer/swerve/cache"

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

func CheckProxy(c *cache.Cache, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// exists, err := certExists(r.Host, c)
		// if err != nil {
		//	log.Error(err)
		// }

		//if r.TLS != nil {
		//	log.Debugf("call next handler, because SSL Cert for Domain %s is in autocert cache", r.Host)
		//	next.ServeHTTP(w, r)
		// }

		host := r.Host
		headerHost := r.Header.Get("X-SWERVE-Forwarded-Host")
		if headerHost != "" {
			host = headerHost
		}

		log.Debugf("call CheckProxy, next lookup for existing certorder for domain %s", host)
		order, _ := checkCertOrder(host, c)
		log.Debugf("CheckProxy: %s", host)

		if order.Hostname != "" {
			target := fmt.Sprintf("http://%s:8080", order.Hostname)
			u, _ := url.Parse(target)
			log.Infof(`CALL REVERSE proxy, forward req to pod %s`, order.Hostname)
			proxy := httputil.NewSingleHostReverseProxy(u)

			r.URL.Host = host
			r.URL.Scheme = "http"
			r.Header.Set("X-SWERVE-Forwarded-Host", host)
			r.Host = host

			proxy.ServeHTTP(w, r)
			r.Context().Done()
			log.Debug("Serve Proxy DONE")
		} else {
			log.Debug("CheckProxy call normal next.Handler")
			next.ServeHTTP(w, r)
		}
	})
}

func certExists(domain string, c *cache.Cache) (bool, error) {
	var (
		data []byte
		err  error
	)
	done := make(chan struct{})

	go func() {
		defer close(done)
		data, err = c.DB.GetCacheEntry(domain)
	}()

	select {
	case <-done:
	case <-time.After(time.Minute):
		return false, fmt.Errorf("checkCertOrder running in timeout")
	}

	if err == nil && data != nil && len(data) > 0 {
		return true, nil
	}

	if err != nil {
		return false, err
	}
	return false, fmt.Errorf("certExists false ...")
}

func checkCertOrder(domain string, c *cache.Cache) (database.CertOrder, error) {
	ip := os.Getenv("SWERVE_POD_IP")
	order, err := c.DB.GetCertOrderEntry(domain)
	if err != nil {
		return order, err
	}
	if order.Domain != "" {
		log.Infof("[name] @ cache GET found CERT ORDER for %s", order.Domain)
	}
	if order.Hostname != ip {
		log.Infof("[name] @ cache IP (local %s) IS NOT EQUAL (order %s) proxy Forward is needed!", ip, order.Hostname)
	}

	return order, nil

}

// LoggingMiddleware logs request infos
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := logWriter{ResponseWriter: w}
		next.ServeHTTP(&sw, r)
		diff := time.Now().Sub(start)
		ts := time.Now().Format("02/Jan/2006 03:04:05")
		log.Infof(`ts="%s" method="%s" proto="%s" code="%d" host="%s" path="%s" qstring="%v" took="%.03fms" ua="%s"`, ts, r.Method, r.Proto, sw.status, r.Host, r.URL.Path, r.URL.RawQuery, float64(diff.Microseconds())/1000, r.UserAgent())
	})
}
