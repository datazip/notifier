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

// ExceptionHandlerMiddleware is a handlers which notfies on slack if any exceptions are encounters in
// http requests such as status code >= 400
func ExceptionHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create custom router
		rw := NewResponseWriter(w)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logrus.Error(err)
			NotifyError("Exception Handler Middleware Error", "failed to read request body", err.Error())
			return
		}

		// setup recovery
		defer func() {
			err := recover()
			if err != nil {
				// capture stacks trace
				stackTrace := string(debug.Stack())
				logrus.Error(err) // May be log this error?
				// print stack trace as well
				fmt.Println(stackTrace)
				w.WriteHeader(http.StatusInternalServerError)
				NotifyError("Exception Handler Middleware Recovery", "Check stack trace above", fmt.Sprintf("%v\n%s", err, stackTrace), "Request", string(body))
			}
		}()

		// clone request and send it ahead
		cloneRequest := r.Clone(r.Context())
		cloneRequest.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(rw, cloneRequest)

		description := fmt.Sprintf("failed request with StatusCode: %d", rw.statusCode)

		if rw.statusCode >= 500 {
			NotifyError("Exception Handler Middleware Error", description, "", "url", r.RequestURI, "Response", string(rw.response), "Request", string(body))
		} else if rw.statusCode >= 400 {
			NotifyWarn("Exception Handler Middleware Warn", description, fmt.Sprintf("Response: %s", string(rw.response)))
		}
	})
}
