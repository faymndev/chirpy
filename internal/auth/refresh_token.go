package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/faymndev/chirpy/internal/database"
)

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}

func IsValidRefreshToken(token database.RefreshToken) bool {
	return token.RevokedAt.Valid && token.ExpiresAt.Compare(time.Now()) == 1
}
