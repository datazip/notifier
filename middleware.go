package notifier

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
	response   []byte
}

func NewResponseWriter(w http.ResponseWriter) *customResponseWriter {
	return &customResponseWriter{w, http.StatusOK, nil}
}

func (rw *customResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *customResponseWriter) Write(content []byte) (int, error) {
	rw.response = content
	return rw.ResponseWriter.Write(content)
}

// ExceptionHandlerMiddlware is a handlers which notfies on slack if any exceptions are encounters in
// http requests such as status code >= 400
func ExceptionHandlerMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create custom router
		rw := NewResponseWriter(w)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logrus.Error(err)
			NotifyError("Exception Handler Middleware Error", "failed to read request body", err.Error())
			return
		}

		// clone request and send it ahead
		cloneRequest := r.Clone(r.Context())
		cloneRequest.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(rw, cloneRequest)

		description := fmt.Sprintf("StatusCode: %d\nRequest body: %s", rw.statusCode, string(body))

		if rw.statusCode >= 500 {
			if debug.Stack() != nil && string(debug.Stack()) != "" {
				NotifyError("Exception Handler Middleware Error", description, fmt.Sprintf("Response: %s\nStack Trace: \n%s", string(rw.response), string(debug.Stack())))
			} else {
				NotifyError("Exception Handler Middleware Error", description, fmt.Sprintf("Response: %s", string(rw.response)))
			}
		} else if rw.statusCode >= 400 {
			NotifyWarn("Exception Handler Middleware Warn", description, fmt.Sprintf("Response: %s", string(rw.response)))
		}
	})
}
