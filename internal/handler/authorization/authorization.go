package authorization

import (
	"net/http"
	"strings"
)

func ValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if bearerToken[0] != "foo" || bearerToken[1] != "bar" {
				http.Error(w, http.StatusText(403), 403)
				return
			}
			next.ServeHTTP(w, req)
		} else {
			http.Error(w, http.StatusText(401), 401)
		}
	})
}
