package handlers

import (
	"net/http"
	"strings"

	jwtauth0 "github.com/auth0/go-jwt-middleware"
	jwtgo "github.com/dgrijalva/jwt-go"
)

// JwtValidator - Negroni-compatible validator
//              - TODO Expose configuration of inner auth0 middleware
type JwtValidator struct {
	handler      *jwtauth0.JWTMiddleware
	ignoredPaths []string
}

// NewJwtValidator -
func NewJwtValidator(ignoredPaths []string, signingKey []byte) *JwtValidator {
	return &JwtValidator{
		jwtauth0.New(jwtauth0.Options{
			ValidationKeyGetter: func(token *jwtgo.Token) (interface{}, error) {
				return signingKey, nil
			},
			SigningMethod: jwtgo.SigningMethodHS256,
		}),
		ignoredPaths,
	}
}

func (v *JwtValidator) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if v.pathIgnored(r.RequestURI) {
		next(rw, r)
	} else {
		err := v.handler.CheckJWT(rw, r)
		// If there was an error, do not call next.
		if err == nil && next != nil {
			next(rw, r)
		}
	}
}

func (v *JwtValidator) pathIgnored(path string) bool {
	for _, s := range v.ignoredPaths {
		if strings.HasPrefix(path, s) {
			return true
		}
	}
	return false
}
