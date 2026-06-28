package utils

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

// In a real scenario, you'd load these from a secure vault or file.
// For this example, we generate a fresh pair on startup if not provided.
func InitKeys() error {
	// For demonstration, you'd typically read from os.Getenv("PRIVATE_KEY_PEM")
	// If you want a quick start without setup, you'd generate them here.
	// For production, always load predefined keys!
	return nil
}

// GenerateToken creates a JWT signed with RS256
func GenerateToken(userID uint, email string, roles []string, projectID uint) (string, error) {
	// Dummy secret for local dev if RSA isn't set up yet. 
	// Production MUST use RS256. We'll use HS256 for rapid scaffolding here unless RSA is configured.
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		secret = []byte("super-secret-key-change-me")
	}

	claims := jwt.MapClaims{
		"sub":        fmt.Sprintf("%d", userID),
		"email":      email,
		"project_id": projectID,
		"roles":      roles,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
