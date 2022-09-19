package main

import (
	"io"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/frozenbonito/unisrv"
)

func TestConfig(t *testing.T) {
	const dynamicPort = 50000

	cases := []struct {
		name        string
		cfg         *config
		validateErr string
		normalized  *config
		addr        string
		url         string
		opts        *unisrv.Options
	}{
		{
			name: "minimum",
			cfg: &config{
				host: "localhost",
			},
			normalized: &config{
				dir:  ".",
				host: "localhost",
				base: "/",
			},
			addr: "localhost:0",
			url:  "http://localhost:50000/",
			opts: &unisrv.Options{
				Base:    "/",
				NoCache: true,
			},
		},
		{
			name: "full",
			cfg: &config{
				dir:            "dir",
				host:           "localhost",
				port:           8080,
				base:           "/base/",
				readTimeout:    10,
				writeTimeout:   15,
				disableNoCache: true,
			},
			normalized: &config{
				dir:            "dir",
				host:           "localhost",
				port:           8080,
				base:           "/base/",
				readTimeout:    10,
				writeTimeout:   15,
				disableNoCache: true,
			},
			addr: "localhost:8080",
			url:  "http://localhost:8080/base/",
			opts: &unisrv.Options{
				Base:    "/base/",
				NoCache: false,
			},
		},
		{
			name: "base has no slash prefix",
			cfg: &config{
				dir:  ".",
				host: "localhost",
				base: "base/",
			},
			normalized: &config{
				dir:  ".",
				host: "localhost",
				base: "/base/",
			},
			addr: "localhost:0",
			url:  "http://localhost:50000/base/",
			opts: &unisrv.Options{
				Base:    "/base/",
				NoCache: true,
			},
		},
		{
			name: "base has no slash suffix",
			cfg: &config{
				dir:  ".",
				host: "localhost",
				base: "/base",
			},
			normalized: &config{
				dir:  ".",
				host: "localhost",
				base: "/base/",
			},
			addr: "localhost:0",
			url:  "http://localhost:50000/base/",
			opts: &unisrv.Options{
				Base:    "/base/",
				NoCache: true,
			},
		},
		{
			name: "base has no slash prefix and suffix",
			cfg: &config{
				dir:  ".",
				host: "localhost",
				base: "base",
			},
			normalized: &config{
				dir:  ".",
				host: "localhost",
				base: "/base/",
			},
			addr: "localhost:0",
			url:  "http://localhost:50000/base/",
			opts: &unisrv.Options{
				Base:    "/base/",
				NoCache: true,
			},
		},
		{
			name:        "host is empty",
			cfg:         &config{},
			validateErr: "host is required",
		},
		{
			name: "port is below minimum value",
			cfg: &config{
				host: "localhost",
				port: -1,
			},
			validateErr: "invalid port",
		},
		{
			name: "port exceeds maximum value",
			cfg: &config{
				host: "localhost",
				port: 65536,
			},
			validateErr: "invalid port",
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			tt.Run("validate", func(ttt *testing.T) {
				err := v.cfg.validate()
				switch {
				case err != nil && err.Error() != v.validateErr:
					ttt.Errorf("expected %q, but got %q", v.validateErr, err.Error())
				case err == nil && v.validateErr != "":
					ttt.Errorf("unexpected success")
				default:
					// nop
				}

				if ttt.Failed() {
					tt.Fail()
				}
			})

			if v.validateErr != "" || tt.Failed() {
				return
			}

			tt.Run("normalize", func(ttt *testing.T) {
				v.cfg.normalize()

				if !reflect.DeepEqual(v.cfg, v.normalized) {
					ttt.Errorf("expected %#v, but got %#v", v.normalized, v.cfg)
				}
			})

			tt.Run("addr", func(ttt *testing.T) {
				addr := v.cfg.addr()
				if addr != v.addr {
					ttt.Errorf("expected %q, but got %q", v.addr, addr)
				}
			})

			tt.Run("url", func(ttt *testing.T) {
				port := v.cfg.port
				if port == 0 {
					port = dynamicPort
				}

				url := v.cfg.url(port)
				if url != v.url {
					ttt.Errorf("expected %q, but got %q", v.url, url)
				}
			})

			tt.Run("serverOptions", func(ttt *testing.T) {
				opts := v.cfg.serverOptions()
				if !reflect.DeepEqual(opts, v.opts) {
					ttt.Errorf("expected %#v, but got %#v", v.opts, opts)
				}
			})
		})
	}
}

func TestParseCommandLineArgs(t *testing.T) {
	// Unset environment variables for test.
	envKeys := []string{
		"UNISRV_HOST",
		"UNISRV_PORT",
		"UNISRV_BASE",
		"UNISRV_READ_TIMEOUT",
		"UNISRV_WRITE_TIMEOUT",
		"UNISRV_DISABLE_NO_CACHE",
	}
	for _, key := range envKeys {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}

	cases := []struct {
		name         string
		env          map[string]string
		args         []string
		cfg          *config
		printVersion bool
		failed       bool
	}{
		{
			name: "default",
			args: []string{},
			cfg: &config{
				host:         "localhost",
				port:         defaultPort,
				readTimeout:  defaultReadTimeout,
				writeTimeout: defaultWriteTimeout,
			},
		},
		{
			name: "env vars",
			env: map[string]string{
				"UNISRV_HOST":             "127.0.0.1",
				"UNISRV_PORT":             "8080",
				"UNISRV_BASE":             "/base1/",
				"UNISRV_READ_TIMEOUT":     "10",
				"UNISRV_WRITE_TIMEOUT":    "15",
				"UNISRV_DISABLE_NO_CACHE": "true",
			},
			args: []string{},
			cfg: &config{
				host:           "127.0.0.1",
				port:           8080,
				base:           "/base1/",
				readTimeout:    10,
				writeTimeout:   15,
				disableNoCache: true,
			},
		},
		{
			name: "args",
			env: map[string]string{
				"UNISRV_HOST":             "127.0.0.1",
				"UNISRV_PORT":             "8080",
				"UNISRV_BASE":             "/base1/",
				"UNISRV_READ_TIMEOUT":     "10",
				"UNISRV_WRITE_TIMEOUT":    "15",
				"UNISRV_DISABLE_NO_CACHE": "true",
			},
			args: []string{
				"-host", "1.1.1.1",
				"-port", "9000",
				"-base", "/base2/",
				"-read-timeout", "20",
				"-write-timeout", "25",
				"-disable-no-cache=false",
				"dir",
			},
			cfg: &config{
				dir:            "dir",
				host:           "1.1.1.1",
				port:           9000,
				base:           "/base2/",
				readTimeout:    20,
				writeTimeout:   25,
				disableNoCache: false,
			},
		},
		{
			name: "print version",
			args: []string{
				"-version",
			},
			printVersion: true,
		},
		{
			name: "invalid args",
			args: []string{
				"-port", "http",
			},
			failed: true,
		},
		{
			name: "too many args",
			args: []string{
				"dir1", "dir2",
			},
			failed: true,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			for key, value := range v.env {
				tt.Setenv(key, value)
			}

			cfg, printVersion, err := parseCommandLineArgs(v.args)

			switch {
			case err != nil && !v.failed:
				tt.Errorf("unexpected error")
			case err == nil && v.failed:
				tt.Errorf("unexpected success")
			default:
				// nop
			}

			if v.failed || tt.Failed() {
				return
			}

			tt.Run("printVersion", func(ttt *testing.T) {
				if printVersion != v.printVersion {
					ttt.Errorf("expected %v, but got %v", v.printVersion, printVersion)
				}
			})

			if v.printVersion {
				return
			}

			tt.Run("config", func(ttt *testing.T) {
				if !reflect.DeepEqual(cfg, v.cfg) {
					ttt.Errorf("expected %#v, but got %#v", v.cfg, cfg)
				}
			})
		})
	}
}

func TestNewServer(t *testing.T) {
	cases := []struct {
		name string
		cfg  *config
	}{
		{
			name: "default",
			cfg: &config{
				dir:  "testdata",
				host: "localhost",
				port: 5000,
				base: "/",
			},
		},
		{
			name: "with base path",
			cfg: &config{
				dir:  "testdata",
				host: "localhost",
				port: 5000,
				base: "/base/",
			},
		},
		{
			name: "with timeout",
			cfg: &config{
				dir:          "testdata",
				host:         "localhost",
				port:         5000,
				base:         "/",
				readTimeout:  5,
				writeTimeout: 10,
			},
		},
	}

	client := http.DefaultClient

	for _, v := range cases {
		t.Run(v.name, func(tt *testing.T) {
			srv := newServer(v.cfg)
			defer srv.Close()

			tt.Run("read timeout", func(ttt *testing.T) {
				expected := time.Duration(v.cfg.readTimeout) * time.Second
				if srv.ReadTimeout != expected {
					ttt.Errorf("expected %v, but got %v", expected, srv.ReadTimeout)
				}
			})

			tt.Run("write timeout", func(ttt *testing.T) {
				expected := time.Duration(v.cfg.writeTimeout) * time.Second
				if srv.WriteTimeout != expected {
					ttt.Errorf("expected %v, but got %v", expected, srv.WriteTimeout)
				}
			})

			listener, err := net.Listen("tcp", v.cfg.addr())
			if err != nil {
				tt.Fatalf("listen failed: %+v", err)
			}
			defer listener.Close()

			go srv.Serve(listener) //nolint:errcheck
			defer srv.Close()

			req, err := http.NewRequest(http.MethodGet, v.cfg.url(v.cfg.port), nil)
			if err != nil {
				tt.Fatalf("failed to create request: %+v", err)
			}

			resp, err := client.Do(req)
			if err != nil {
				tt.Fatalf("request failed: %+v", err)
			}
			defer resp.Body.Close()

			tt.Run("status code", func(ttt *testing.T) {
				if resp.StatusCode != http.StatusOK {
					ttt.Errorf("expected %d, but got %d", http.StatusOK, resp.StatusCode)
				}
			})

			tt.Run("body", func(ttt *testing.T) {
				buf := &strings.Builder{}
				if _, err := io.Copy(buf, resp.Body); err != nil {
					ttt.Fatalf("copy failed: %+v", err)
				}

				body := buf.String()
				expected := "testdata\n"
				if body != expected {
					ttt.Errorf("expected %q, but got %q", expected, body)
				}
			})
		})
	}
}
