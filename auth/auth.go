package auth

import (
	"context"
	"e-market/model"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type key string
type MyClaims struct {
	model.User
	jwt.RegisteredClaims
}

var (
	MySigningKey     = []byte(os.Getenv("MY_SECRET_KEY"))
	Metadata     key = "katseye"
)

func CreateSS(User model.User) (string, error) {
	claims := MyClaims{
		User: User,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(28 * time.Hour)),
			Issuer:    "Me!!!!!!!!!!!!",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(MySigningKey)
	if err != nil {
		return "", nil
	}
	return ss, nil
}

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
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), Metadata, claims.User)

			next(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized: token invalid", http.StatusUnauthorized)
			return
		}
	}
}
