package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// RestServer encapsulates everything needed to run the API
type RestServer struct {
	// Add other internal dependencies as needed
	router *httprouter.Router
}

// New creates and initializes the REST server and its resources
func New() *RestServer {
	s := &RestServer{
		router: httprouter.New(),
	}
	s.router.PanicHandler = func(rw http.ResponseWriter, r *http.Request, v interface{}) {
		log.Printf("%s -- Internal Server Error: %v \n", r.RequestURI, v)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	s.router.GET("/v1/resource", notImplemented())
	s.router.POST("/v1/resource", create())
	s.router.PUT("/v1/resource/:id", notImplemented())
	s.router.DELETE("/v1/resource/:id", notImplemented())
	s.router.GET("/healthz", health)
	return s
}

// Serve starts the server on the specified port, optionally
// performing request logging...
func (s *RestServer) Serve(httpPort int, logRequests bool) error {
	log.Printf("HTTP listening on %d", httpPort)
	var h http.Handler
	if logRequests {
		h = NewRequestLogger(s.router)
	} else {
		h = s.router
	}
	return http.ListenAndServe(fmt.Sprintf(":%d", httpPort), h)
}

func create() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, "TODO")
	}
}

func notImplemented() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")
}
