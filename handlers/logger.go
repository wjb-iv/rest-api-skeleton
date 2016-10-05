package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/urfave/negroni"
)

const apacheFormatPattern = "%s - - [%s] \"%s %d %d\" %f\n"

// ApacheLogger is a middleware handler that logs the request in Apache standard format
type ApacheLogger struct {
	// Just a target for the ServeHttp method
}

// NewApacheLogger - returns a logger
func NewApacheLogger() *ApacheLogger {
	return &ApacheLogger{}
}

func (al *ApacheLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	clientIP := r.RemoteAddr
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}
	// Call the next handler with timings...
	startTime := time.Now()
	next.ServeHTTP(rw, r)
	finishTime := time.Now()
	// Use some Negroni goodness to obtain info from response
	resp := rw.(negroni.ResponseWriter)

	log.Printf(apacheFormatPattern,
		clientIP,
		finishTime.UTC().Format("02/Jan/2006 03:04:05"),
		fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto),
		resp.Status(),
		resp.Size(),
		finishTime.Sub(startTime).Seconds())
}
