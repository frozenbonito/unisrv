package unisrv_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/frozenbonito/unisrv"
)

func TestNewHandler(t *testing.T) {
	cases := []struct {
		name    string
		opts    *unisrv.Options
		noCache bool
	}{
		{
			name: "nil options",
			opts: nil,
		},
		{
			name: "with base option",
			opts: &unisrv.Options{
				Base: "/",
			},
		},
		{
			name: "with base option",
			opts: &unisrv.Options{
				Base: "/base/",
			},
		},
		{
			name: "with no cache option",
			opts: &unisrv.Options{
				NoCache: true,
			},
			noCache: true,
		},
	}

	targets := []struct {
		path            string
		statusCode      int
		contentEncoding string
		contentType     string
	}{
		{
			path:        "/",
			statusCode:  http.StatusOK,
			contentType: "text/html; charset=utf-8",
		},
		{
			path:       "/index.html",
			statusCode: http.StatusMovedPermanently,
		},
		{
			path:            "/Build/Build.data.br",
			statusCode:      http.StatusOK,
			contentEncoding: "br",
			contentType:     "application/octet-stream",
		},
		{
			path:            "/Build/Build.symbols.json.br",
			statusCode:      http.StatusOK,
			contentEncoding: "br",
			contentType:     "application/octet-stream",
		},
		{
			path:            "/Build/Build.framework.js.br",
			statusCode:      http.StatusOK,
			contentEncoding: "br",
			contentType:     "application/javascript",
		},
		{
			path:            "/Build/Build.wasm.br",
			statusCode:      http.StatusOK,
			contentEncoding: "br",
			contentType:     "application/wasm",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			h := unisrv.NewHandler("testdata", v.opts)

			base := "/"
			if v.opts != nil && v.opts.Base != "" {
				base = v.opts.Base
			}

			for _, target := range targets {
				tt.Run(target.path, func(ttt *testing.T) {
					targetPath := path.Join(base, target.path)
					if strings.HasSuffix(target.path, "/") {
						targetPath += "/"
					}

					r := httptest.NewRequest(http.MethodGet, targetPath, nil)
					w := httptest.NewRecorder()

					h.ServeHTTP(w, r)

					resp := w.Result()
					defer resp.Body.Close()

					ttt.Run("status code", func(tttt *testing.T) {
						if resp.StatusCode != target.statusCode {
							tttt.Errorf("expected %d, but got %d", target.statusCode, resp.StatusCode)
						}
					})

					ttt.Run("Content-Encoding header", func(tttt *testing.T) {
						contentEncoding := resp.Header.Get("Content-Encoding")
						if contentEncoding != target.contentEncoding {
							tttt.Errorf("expected %q, but got %q", target.contentEncoding, contentEncoding)
						}
					})

					ttt.Run("Content-Type header", func(tttt *testing.T) {
						contentType := resp.Header.Get("Content-Type")
						if contentType != target.contentType {
							tttt.Errorf("expected %q, but got %q", target.contentType, contentType)
						}
					})

					ttt.Run("Cache-Control header", func(tttt *testing.T) {
						cacheControl := resp.Header.Get("Cache-Control")

						expected := ""
						if v.noCache {
							expected = "no-cache"
						}

						if cacheControl != expected {
							tttt.Errorf("expected %q, but got %q", expected, cacheControl)
						}
					})
				})
			}
		})
	}
}

func TestUnityMiddleware(t *testing.T) {
	cases := []struct {
		path            string
		contentEncoding string
		contentType     string
	}{
		{
			path:            "/Build/Build.data.br",
			contentEncoding: "br",
			contentType:     "application/octet-stream",
		},
		{
			path:            "/Build/Build.symbols.json.br",
			contentEncoding: "br",
			contentType:     "application/octet-stream",
		},
		{
			path:            "/Build/Build.framework.js.br",
			contentEncoding: "br",
			contentType:     "application/javascript",
		},
		{
			path:            "/Build/Build.wasm.br",
			contentEncoding: "br",
			contentType:     "application/wasm",
		},
		{
			path:            "/Build/Build.data.gz",
			contentEncoding: "gzip",
			contentType:     "application/octet-stream",
		},
		{
			path:            "/Build/Build.symbols.json.gz",
			contentEncoding: "gzip",
			contentType:     "application/octet-stream",
		},
		{
			path:            "/Build/Build.framework.js.gz",
			contentEncoding: "gzip",
			contentType:     "application/javascript",
		},
		{
			path:            "/Build/Build.wasm.gz",
			contentEncoding: "gzip",
			contentType:     "application/wasm",
		},
	}

	h := unisrv.UnityMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "ok")
	}))

	for _, v := range cases {
		t.Run(v.path, func(tt *testing.T) {
			r := httptest.NewRequest(http.MethodGet, v.path, nil)
			w := httptest.NewRecorder()

			h.ServeHTTP(w, r)

			resp := w.Result()
			defer resp.Body.Close()

			tt.Run("status code", func(ttt *testing.T) {
				if resp.StatusCode != http.StatusOK {
					ttt.Errorf("expected %d, but got %d", http.StatusOK, resp.StatusCode)
				}
			})

			tt.Run("Content-Encoding header", func(ttt *testing.T) {
				contentEncoding := resp.Header.Get("Content-Encoding")
				if contentEncoding != v.contentEncoding {
					ttt.Errorf("expected %q, but got %q", v.contentEncoding, contentEncoding)
				}
			})

			tt.Run("Content-Type header", func(ttt *testing.T) {
				contentType := resp.Header.Get("Content-Type")
				if contentType != v.contentType {
					ttt.Errorf("expected %q, but got %q", v.contentType, contentType)
				}
			})
		})
	}
}
