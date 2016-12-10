package rest

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const apacheFormatPattern = "%s - - [%s] \"%s %d %d\" %f\n"

type RequestLogger struct {
	handler http.Handler
}

func NewRequestLogger(handler http.Handler) *RequestLogger {
	return &RequestLogger{handler: handler}
}

func (rl *RequestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.RequestURI, "healthz") {
		rl.handler.ServeHTTP(w, r)
		return
	}
	clientIP := r.RemoteAddr
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}
	// Wrapping the response writer lets us gather some extra data
	rww := &responseWriterWrapper{wrapped: w, status: http.StatusOK}
	// Call the inner http handler with the wrapped writer and timing
	startTime := time.Now()
	rl.handler.ServeHTTP(rww, r)
	finishTime := time.Now()

	log.Printf(apacheFormatPattern,
		clientIP,
		finishTime.UTC().Format("02/Jan/2006 03:04:05"),
		fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto),
		rww.Status(),
		rww.Size(),
		finishTime.Sub(startTime).Seconds())
}

type responseWriterWrapper struct {
	wrapped http.ResponseWriter
	status  int
	size    int
}

func (rw *responseWriterWrapper) Header() http.Header {
	return rw.wrapped.Header()
}

func (rw *responseWriterWrapper) Write(buf []byte) (int, error) {
	size, err := rw.wrapped.Write(buf)
	rw.size += size
	return size, err
}

func (rw *responseWriterWrapper) WriteHeader(st int) {
	rw.status = st
	rw.wrapped.WriteHeader(st)
}

func (rw *responseWriterWrapper) Status() int {
	return rw.status
}

func (rw *responseWriterWrapper) Size() int {
	return rw.size
}
