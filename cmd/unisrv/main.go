package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/frozenbonito/unisrv"
	"github.com/frozenbonito/unisrv/internal/middleware"
)

const (
	defaultPort         = 5000
	defaultReadTimeout  = 5
	defaultWriteTimeout = 5
)

var version = "dev"

// config is cli config.
type config struct {
	dir            string
	host           string
	port           int
	base           string
	readTimeout    int
	writeTimeout   int
	disableNoCache bool
}

// validate reports whether the config is valid.
func (s *config) validate() error {
	if s.host == "" {
		return errors.New("host is required")
	}
	if s.port < 0 || s.port > 65535 {
		return errors.New("invalid port")
	}
	return nil
}

// normalize normalizes the config.
func (s *config) normalize() {
	if s.dir == "" {
		s.dir = "."
	}

	if s.base == "" {
		s.base = "/"
	} else {
		if !strings.HasPrefix(s.base, "/") {
			s.base = "/" + s.base
		}
		if !strings.HasSuffix(s.base, "/") {
			s.base += "/"
		}
	}
}

// addr returns a TCP network address for a server.
func (s *config) addr() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

// url returns a URL of Unity application.
func (s *config) url(port int) string {
	return fmt.Sprintf("http://%s:%d%s", s.host, port, s.base)
}

// serverOptions returns options for unisrv handler.
func (s *config) serverOptions() *unisrv.Options {
	return &unisrv.Options{
		Base:    s.base,
		NoCache: !s.disableNoCache,
	}
}

func main() {
	cfg, printVersion, err := parseCommandLineArgs(os.Args[1:])
	if err != nil {
		os.Exit(2) //nolint:mnd
	}

	if printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	ctx := context.Background()

	if err := run(ctx, cfg); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// parseCommandLineArgs parses given arguments and then returns config and whether the version should be printed.
func parseCommandLineArgs(args []string) (cfg *config, printVersion bool, err error) {
	cfg = &config{}

	fs := flag.NewFlagSet("unisrv", flag.ContinueOnError)

	fs.StringVar(&cfg.host, "host", "localhost", "hostname")
	fs.IntVar(&cfg.port, "port", defaultPort, "port number")
	fs.StringVar(&cfg.base, "base", "", "base path")
	fs.IntVar(&cfg.readTimeout, "read-timeout", defaultReadTimeout, "maximum duration for reading request in seconds")
	fs.IntVar(&cfg.writeTimeout, "write-timeout", defaultWriteTimeout, "maximum duration for writing response in seconds")
	fs.BoolVar(&cfg.disableNoCache, "disable-no-cache", false, "disable setting 'Cache-Control: no-cache' header")
	fs.BoolVar(&printVersion, "version", false, "print version")

	fs.VisitAll(func(f *flag.Flag) {
		if f.Name == "version" {
			return
		}

		key := fmt.Sprintf("UNISRV_%s", strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_")))
		if s := os.Getenv(key); s != "" {
			f.Value.Set(s) //nolint:errcheck
		}
	})

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: unisrv [flags] [path]\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, false, fmt.Errorf("parse flags: %w", err)
	}

	nonFlagArgs := fs.Args()
	if len(nonFlagArgs) > 1 {
		fs.Usage()
		return nil, false, errors.New("too many arguments")
	}

	if len(nonFlagArgs) == 1 {
		cfg.dir = nonFlagArgs[0]
	}

	return cfg, printVersion, nil
}

// run starts server.
func run(ctx context.Context, cfg *config) error {
	if err := cfg.validate(); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	cfg.normalize()

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	srv := newServer(cfg)

	listener, err := net.Listen("tcp", cfg.addr())
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	errChan := make(chan error)

	fmt.Printf("server running at: %s\n", cfg.url(port))
	go func() {
		if err := srv.Serve(listener); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				close(errChan)
				return
			}
			errChan <- fmt.Errorf("serve: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		if err := srv.Shutdown(context.Background()); err != nil {
			fmt.Fprintln(os.Stderr, "filed to shutdown server:", err)
		}
	}

	return <-errChan
}

// newServer creates a new server.
func newServer(cfg *config) *http.Server {
	mux := http.NewServeMux()

	h := unisrv.NewHandler(cfg.dir, cfg.serverOptions())
	h = middleware.RequestLogger(h)
	mux.Handle(cfg.base, h)

	return &http.Server{
		Handler:      mux,
		ReadTimeout:  time.Duration(cfg.readTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.writeTimeout) * time.Second,
	}
}
