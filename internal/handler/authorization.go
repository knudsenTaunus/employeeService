package handler

import (
	"net/http"
	"strings"
)

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if bearerToken[0] != "foo" || bearerToken[1] != "bar" {
				http.Error(w, http.StatusText(403),403)
				return
			}
			next(w,req)
		} else {
			http.Error(w, http.StatusText(401),401)
		}
	}
}
