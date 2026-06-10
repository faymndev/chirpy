package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func getSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func MakeJWT(userID uuid.UUID, expiresIn time.Duration) (string, error) {
	issuedAt := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(issuedAt),
		ExpiresAt: jwt.NewNumericDate(issuedAt.Add(expiresIn)),
		Subject:   userID.String(),
	})

	return token.SignedString(getSecret())
}

// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#ParseWithClaims
func VerifyJWT(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return getSecret(), nil
	})

	if err != nil {
		return uuid.UUID{}, err
	} else if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return uuid.Parse(claims.Subject)
	}

	return uuid.UUID{}, errors.New("Failed to parse claims.")
}
