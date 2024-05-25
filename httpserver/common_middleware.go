package httpserver

import (
	"context"
	"net/http"
)

const request_context_key = ContextKey("request")
const response_writer_context_key = ContextKey("response_writer")

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
