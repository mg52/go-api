package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/mg52/go-api/domain"
	"github.com/mg52/go-api/helper"
	"log"
	"net/http"
	"os"
	"strings"
)

func RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing RequireAuthentication Middleware")
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			//w.Write([]byte("Unauthorized"))
			return
		}

		jwtToken := strings.Split(token, "Bearer ")
		// Initialize a new instance of `Claims`
		claims := &domain.Claims{}
		tkn, err := jwt.ParseWithClaims(jwtToken[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil {
			helper.Unauthorized(w, r, err)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
		log.Print("Executing RequireAuthentication Middleware again")
	})
}
