// Package unisrv implements http helpers for serving Unity WebGL application.
package unisrv

import (
	"net/http"
	"path"
	"strings"
)

// NewHandler returns a handler that serves Unity application.
func NewHandler(dir string, opts *Options) http.Handler {
	if opts == nil {
		opts = &Options{}
	}

	h := http.FileServer(http.Dir(dir))

	if opts.Base != "" && opts.Base != "/" {
		h = http.StripPrefix(opts.Base, h)
	}

	h = UnityMiddleware(h)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if opts.NoCache {
			w.Header().Set("Cache-Control", "no-cache")
		}
		h.ServeHTTP(w, r)
	})
}

// UnityMiddleware is a middleware for serving Unity application.
func UnityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncoding, contentType := contentHeaders(r)
		if contentEncoding != "" {
			w.Header().Set("Content-Encoding", contentEncoding)
		}
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		next.ServeHTTP(w, r)
	})
}

// contentHeaders returns values of `Content-Encoding` and `Content-Type` response headers
// for serving Unity application.
func contentHeaders(r *http.Request) (contentEncoding, contentType string) {
	ext := path.Ext(r.URL.Path)
	switch ext {
	case ".br":
		contentEncoding = "br"
	case ".gz":
		contentEncoding = "gzip"
	default:
		return
	}

	original := strings.TrimSuffix(r.URL.Path, ext)
	switch {
	case strings.HasSuffix(original, ".data"), strings.HasSuffix(original, ".symbols.json"):
		contentType = "application/octet-stream"
	case strings.HasSuffix(original, ".js"):
		contentType = "application/javascript"
	case strings.HasSuffix(original, ".wasm"):
		contentType = "application/wasm"
	}
	return
}
