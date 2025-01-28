package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

type UserIDKey string

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Токен не найден", http.StatusUnauthorized)
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("denet_secret_key_here"), nil
		})

		if err != nil {
			http.Error(w, "Неправильный токен", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			http.Error(w, "Неправильный токен", http.StatusUnauthorized)
			return
		}

		if token.Method != jwt.SigningMethodHS256 {
			http.Error(w, "Неправильный метод подписи токена", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey("user_id"), claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
