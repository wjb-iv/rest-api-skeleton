package resources

import (
	"encoding/json"
	"log"
	"net/http"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

type temp struct {
	Path string `json:"path"`
	ID   string `json:"id"`
	Msg  string `json:"msg"`
	//Auth string `json:"auth"`
}

// Init -
func Init(r *httprouter.Router) {
	r.GET("/resource", list())
	r.POST("/resource", create())
	r.GET("/resource/:id", get())
	r.PATCH("/resource/:id", update())
	r.DELETE("/resource/:id", remove())
}

func list() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		dumpUser(r)
		res := temp{r.RequestURI, "", "Not Implemented Yet"}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func get() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		idVar := params.ByName("id")
		res := temp{r.RequestURI, idVar, "Not Implemented Yet"}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
	}
}

func create() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		res := temp{r.RequestURI, "", "Not Implemented Yet"}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.Header().Add("Location", "resource/null")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	}
}

func update() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		//idVar := params.ByName("id")
		w.WriteHeader(http.StatusNoContent)
	}
}

func remove() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		//idVar := params.ByName("id")
		w.WriteHeader(http.StatusNoContent)
	}
}

func dumpUser(r *http.Request) {
	if tok := context.Get(r, "user"); tok != nil {
		if token, ok := tok.(*jwtgo.Token); ok {
			if claims, ok := token.Claims.(jwtgo.MapClaims); ok {
				log.Printf("Username: %v; Roles: %v", claims["name"], claims["roles"])
				return
			}
		}
	}
	log.Println("No user found in context")
}
