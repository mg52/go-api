package helper

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mg52/go-api/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

// NewLogger creates new Logger
func NewLogger(environment string) *logrus.Entry {
	l := logrus.New()
	l.Out = os.Stdout
	l.Formatter = &logrus.JSONFormatter{}
	l.Level = logrus.InfoLevel

	return l.WithFields(logrus.Fields{
		"env": environment,
	})
}

func GenerateJWTToken(user domain.User) (string, error) {
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &domain.Claims{
		ID:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", errors.New("GenerateJWTToken Error")
	}
	return tokenString, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
