package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
)

// InitKeys loads the RSA keys from the database, or generates a new pair if they don't exist.
func InitKeys() error {
	var authKey models.AuthKey
	result := database.DB.First(&authKey)

	if result.Error != nil {
		// Generate new RSA key pair
		fmt.Println("Generating new RSA-2048 key pair and saving to database...")
		priv, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return fmt.Errorf("failed to generate RSA key: %v", err)
		}

		privBytes := x509.MarshalPKCS1PrivateKey(priv)
		privPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privBytes,
		})

		pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to marshal public key: %v", err)
		}
		pubPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubBytes,
		})

		authKey = models.AuthKey{
			PrivateKey: string(privPEM),
			PublicKey:  string(pubPEM),
		}
		if err := database.DB.Create(&authKey).Error; err != nil {
			return fmt.Errorf("failed to save keys to database: %v", err)
		}
		fmt.Println("Keys generated and saved to database successfully.")
	} else {
		fmt.Println("Loaded existing RSA keys from database.")
	}

	// Load Private Key
	privBlock, _ := pem.Decode([]byte(authKey.PrivateKey))
	if privBlock == nil {
		return errors.New("failed to decode PEM block containing private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}
	privateKey = priv

	// Load Public Key
	pubBlock, _ := pem.Decode([]byte(authKey.PublicKey))
	if pubBlock == nil {
		return errors.New("failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %v", err)
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return errors.New("not an RSA public key")
	}
	PublicKey = rsaPub

	return nil
}

// GenerateToken creates a JWT signed with RS256
func GenerateToken(userID uint, email string, roles []string, projectID uint) (string, error) {
	if privateKey == nil {
		return "", errors.New("RSA private key not initialized")
	}

	claims := jwt.MapClaims{
		"sub":        fmt.Sprintf("%d", userID),
		"email":      email,
		"project_id": projectID,
		"roles":      roles,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}
