package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func (s *JWTService) ParseAccessToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %s", token.Header["alg"])
		}
		return []byte(s.Salt), nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to parse token: %s", err)
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.Subject == "" || claims.ID == "" || claims.ExpiresAt == nil {
			return nil, fmt.Errorf("Subject claims Subject,ID,ExpiresAt are empty")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("Invalid token")
}
