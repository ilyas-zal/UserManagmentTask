// middlewares содержит функции для работы с токенами и аутентификацией.
package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// Claims содержит данные о пользователе и стандартные данные токена.
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken генерирует токен для пользователя.
func GenerateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
			Issuer:    "denet",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("denet_secret_key_here"))
}

// AuthMiddleware - это функция аутентификации для HTTP-запросов.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					http.Error(w, "Неправильный токен", http.StatusUnauthorized)
					return
				} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
					http.Error(w, "Неправильный токен", http.StatusUnauthorized)
					return
				} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					http.Error(w, "Неправильный токен", http.StatusUnauthorized)
					return
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					http.Error(w, "Токен истек", http.StatusUnauthorized)
					return
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					http.Error(w, "Токен еще не действителен", http.StatusUnauthorized)
					return
				} else {
					http.Error(w, "Неправильный токен", http.StatusUnauthorized)
					return
				}
			}
			http.Error(w, "Неправильный токен", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, "Неправильный токен", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
