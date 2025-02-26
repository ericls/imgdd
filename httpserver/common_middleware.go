package httpserver

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/ericls/imgdd/logging"

	"github.com/felixge/httpsnoop"
	"github.com/rs/zerolog"
)

const request_context_key = ContextKey("request")
const response_writer_context_key = ContextKey("response_writer")

var httpLogger zerolog.Logger

func init() {
	httpLogger = logging.GetLogger("http_logger")
}

func RWContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newContext := context.WithValue(r.Context(), request_context_key, r)
		newContext = context.WithValue(newContext, response_writer_context_key, w)
		r = r.WithContext(newContext)
		next.ServeHTTP(w, r)
	})
}

func GetResponseWriter(c context.Context) http.ResponseWriter {
	v, ok := c.Value(response_writer_context_key).(http.ResponseWriter)
	if v == nil || !ok {
		return nil
	}
	return v
}

func GetRequest(c context.Context) *http.Request {
	v, ok := c.Value(request_context_key).(*http.Request)
	if v == nil || !ok {
		return nil
	}
	return v
}

func IsSecure(r *http.Request) bool {
	if r == nil {
		return false
	}
	if r.TLS != nil {
		return true
	}
	if r.Header.Get("X-Forwarded-Proto") == "https" {
		return true
	}

	return false
}

func IsHttps(c context.Context) bool {
	r := GetRequest(c)
	return IsSecure(r)
}

func GetBaseURL(c context.Context) *url.URL {
	r := GetRequest(c)
	url := &url.URL{
		Scheme: "http",
		Host:   r.Host,
	}
	if IsSecure(r) {
		url.Scheme = "https"
	}
	return url
}

type loggingEntry struct {
	Proto      string
	Method     string
	HostName   string
	Port       string
	URL        string
	StatusCode int
	Size       int64
	RemoteIP   string
	Duration   time.Duration
	Referer    string
}

func (l *loggingEntry) Log() {
	httpLogger.Info().
		Str("remote_ip", l.RemoteIP).
		Str("proto", l.Proto).
		Str("method", l.Method).
		Str("host", l.HostName).
		Str("port", l.Port).
		Str("url", l.URL).
		Str("referer", l.Referer).
		Int("status_code", l.StatusCode).
		Int64("size", l.Size).
		Dur("duration", l.Duration).
		Msg("http request")
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics := httpsnoop.CaptureMetrics(
			next,
			w,
			r,
		)

		remoteHost := ExtractIP(r)

		uri := r.RequestURI

		if r.ProtoMajor == 2 && r.Method == "CONNECT" {
			uri = r.Host
		}
		if uri == "" {
			uri = r.URL.RequestURI()
		}

		host, port, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
			if r.TLS != nil {
				port = "443"
			} else {
				port = "80"
			}
		}

		entry := &loggingEntry{
			Proto:      r.Proto,
			Method:     r.Method,
			URL:        uri,
			StatusCode: metrics.Code,
			Size:       metrics.Written,
			Duration:   metrics.Duration,
			Referer:    r.Referer(),
			RemoteIP:   remoteHost,
			HostName:   host,
			Port:       port,
		}
		entry.Log()
	})
}
