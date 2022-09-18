package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

// RequestLogger is a middleware that prints simple access logs.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{
			ResponseWriter: w,
		}
		defer func() {
			logger.Printf(`"%s %s %s" - %d`, r.Method, r.RequestURI, r.Proto, rw.StatusCode())
		}()

		next.ServeHTTP(rw, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
}

func (s *responseWriter) Write(p []byte) (int, error) {
	if !s.wroteHeader {
		s.WriteHeader(http.StatusOK)
	}

	n, err := s.ResponseWriter.Write(p)
	if err != nil {
		return n, fmt.Errorf("write: %w", err)
	}

	return n, nil
}

func (s *responseWriter) WriteHeader(code int) {
	s.code = code
	s.wroteHeader = true
	s.ResponseWriter.WriteHeader(code)
}

func (s *responseWriter) StatusCode() int {
	return s.code
}
