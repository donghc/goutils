package netutil

import (
	"net/http"

	"github.com/tomasen/realip"
)

func GetRequestIP(r *http.Request) string {
	if r.Header.Get("Proxy-Client-IP") != "" {
		return r.Header.Get("Proxy-Client-IP")
	}
	if r.Header.Get("WL-Proxy-Client-IP") != "" {
		return r.Header.Get("WL-Proxy-Client-IP")
	}
	return realip.FromRequest(r)
}
