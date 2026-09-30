package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	RoleID int    `json:"role_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(
	userID int,
	roleID int,
	role string,
) (string, error) {

	secret := os.Getenv("JWT_SECRET")

	claims := Claims{
		UserID: userID,
		RoleID: roleID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(secret))
}
