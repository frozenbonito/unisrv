package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseWriter(t *testing.T) {
	t.Run("write without writing header", func(tt *testing.T) {
		w := httptest.NewRecorder()
		rw := &responseWriter{
			ResponseWriter: w,
		}

		fmt.Fprint(rw, "ok")

		statusCode := rw.StatusCode()
		if statusCode != http.StatusOK {
			tt.Errorf("expected %d, but got %d", http.StatusOK, statusCode)
		}
	})

	t.Run("write header", func(tt *testing.T) {
		w := httptest.NewRecorder()
		rw := &responseWriter{
			ResponseWriter: w,
		}

		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(rw, "error")

		statusCode := rw.StatusCode()
		if statusCode != http.StatusInternalServerError {
			tt.Errorf("expected %d, but got %d", http.StatusInternalServerError, statusCode)
		}
	})
}
