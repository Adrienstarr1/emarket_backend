package user

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorised missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(t *jwt.Token) (any, error) {
			return MySigningKey, nil
		})
		/*
			so like this parses the claims the compares it with the sigining key
			if the signing key and the parsed token are like the same ok will e tru else you get an error and ok is false
			the &myclaims kinda gives it the parameter its meant to splice the parsed info into
		*/
		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), "user_id", claims.User.Id)
			next(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}
	}
}
