//go:build dev

package front

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Handler() http.HandlerFunc {
	target, _ := url.Parse("http://localhost:5173")
	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}
