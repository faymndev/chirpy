package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestSigning(t *testing.T) {
	userID := uuid.New()

	expiresIn, err := time.ParseDuration("30m")
	if err != nil {
		t.Errorf("Failed to parse duration")
	}

	jwt, err := MakeJWT(userID, expiresIn)
	if err != nil {
		t.Errorf("Failed to make jwt, %v", err)
	}

	id, err := VerifyJWT(jwt)
	if err != nil {
		t.Error("Failed to verify jwt")
	} else if id != userID {
		t.Error("Invalid user ID")
	}
}
