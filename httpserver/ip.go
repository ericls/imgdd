package httpserver

import (
	"net"
	"net/http"
	"strings"
)

func ExtractIPDirect(req *http.Request) string {
	h, _, _ := net.SplitHostPort(req.RemoteAddr)
	return h
}

func ExtractIPXFF(req *http.Request) string {
	xffIps := req.Header["X-Forwarded-For"]
	if len(xffIps) == 0 {
		return ""
	}
	ips := strings.Split(strings.Join(xffIps, ","), ",")
	for i := len(ips) - 1; i >= 0; i-- {
		ips[i] = strings.TrimSpace(ips[i])
		ips[i] = strings.TrimPrefix(ips[i], "[")
		ips[i] = strings.TrimSuffix(ips[i], "]")
		if net.ParseIP(ips[i]) == nil {
			return ""
		}
	}
	return strings.TrimSpace(ips[0])
}

func ExtractIP(req *http.Request) string {
	ip := ExtractIPXFF(req)
	if ip == "" {
		ip = ExtractIPDirect(req)
	}
	return ip
}
