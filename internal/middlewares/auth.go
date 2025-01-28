package middlewares

import (
	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "myapp",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("denet_secret_key_here"))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
