package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wjb-iv/rest-api-skeleton/handlers"
	"github.com/wjb-iv/rest-api-skeleton/resources"

	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
)

func health(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func main() {
	router := httprouter.New()

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(handlers.NewApacheLogger())
	n.Use(handlers.NewJwtValidator([]string{"/login", "/healthz"}, []byte("secret")))

	router.GET("/healthz", health)
	resources.Init(router)

	n.UseHandler(router)

	log.Fatal(http.ListenAndServe(":4100", n))
}
