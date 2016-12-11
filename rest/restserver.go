package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

	// This allows us to test JWT tokens without a separate auth server
	s.router.POST("/token", authenticated(AUTH_BASIC, token()))

	// Our primary API resources are here...
	s.router.GET("/v1/resource", authenticated(AUTH_JWT, list()))
	s.router.GET("/v1/resource/:id", authenticated(AUTH_JWT, get()))
	s.router.POST("/v1/resource", authenticated(AUTH_JWT, create()))
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

// For demo purposes, we implement this as a simple form POST - a real
// implementation would use oauth or similar protocol
func token() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Pull user data out of request...
		grantType := r.PostFormValue("grant_type")
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		if grantType == "" || username == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("User %s has requested a token. \n", username)
		// Create the token
		// TODO - actually validate the user and assign roles
		roles := []string{"ROLE_USER", "ROLE_ADMIN"}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"roles": roles,
			"admin": true,
			"name":  username,
			"exp":   time.Now().Add(time.Hour * 1).Unix(),
		})
		// Sign the token with our secret
		tokenString, _ := token.SignedString(mySigningKey)
		tk := TokenStruct{
			AccessToken:  tokenString,
			ExpiresIn:    3600,
			TokenType:    "bearer",
			RefreshToken: "ToBeImplemented",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&tk)
	}
}

func list() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		dumpUserData(r)
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "TODO")
	}
}

func get() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		dumpUserData(r)
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "TODO")
	}
}

func create() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		dumpUserData(r)
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

func dumpUserData(r *http.Request) {
	username := r.Context().Value("username")
	roles := r.Context().Value("roles")
	log.Printf("Resource accessed by user: %v with roles: %v \n", username, roles)
}

const (
	AUTH_BASIC = 1
	AUTH_JWT   = 2
)

var mySigningKey = []byte("secret")

// TokenStruct is used internally to return json in response to token() request
type TokenStruct struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

func authenticated(auth int, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if auth == AUTH_BASIC {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			usr, pwd, ok := r.BasicAuth()
			// TODO - really validate the username and password
			if !ok || usr != "username" || pwd != "password" {
				http.Error(w, "Not authorized", http.StatusForbidden)
				return
			}
		}
		if auth == AUTH_JWT {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Not authorized", http.StatusForbidden)
				return
			}
			authHeaderParts := strings.Split(authHeader, " ")
			if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusBadRequest)
				return
			}
			// Now we have the token, perform the validation...
			token, err := jwt.Parse(authHeaderParts[1], func(token *jwt.Token) (interface{}, error) { return mySigningKey, nil })
			if token.Valid {
				// The token is valid, so pass along user data in request's context
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					log.Println("Adding JWT claims to request context...")
					newCtx := context.WithValue(r.Context(), "username", claims["name"])
					newCtx = context.WithValue(newCtx, "roles", claims["roles"])
					r = r.WithContext(newCtx)
				}
			} else {
				// Token not valid, this shows some examples of the errors that might happen
				if ve, ok := err.(*jwt.ValidationError); ok {
					if ve.Errors&jwt.ValidationErrorMalformed != 0 {
						// Token is malformed
						log.Printf("Token: %s is malformed. \n", authHeaderParts[1])
					} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
						// Token is either expired or not active yet
						log.Printf("Token: %s is expired or not yet acrive. \n", authHeaderParts[1])
					} else {
						log.Printf("Unexpected error validating token: %v \n", err)
					}
				}
				http.Error(w, "Not authorized", http.StatusForbidden)
				return
			}
		}
		// If we make it here, we execute the wrapped httprouter.Handle
		next(w, r, p)
	}
}
