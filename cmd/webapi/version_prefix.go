package main

import (
	"net/http"
	"strings"
)

func stripPathPrefix(prefix string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil || r.URL == nil {
			next.ServeHTTP(w, r)
			return
		}
		if r.URL.Path == prefix {
			r.URL.Path = "/"
		} else if strings.HasPrefix(r.URL.Path, prefix+"/") {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
		}
		next.ServeHTTP(w, r)
	})
}
