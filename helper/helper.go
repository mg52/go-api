package helper

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mg52/go-api/domain"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	var resp domain.Response
	resp.Msg = "not found"
	jsonBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusNotFound)
	w.Write(jsonBytes)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	var resp domain.Response
	resp.Err = err.Error()
	jsonBytes, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonBytes)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	var resp domain.Response

	w.WriteHeader(http.StatusUnauthorized)
	if err != nil {
		resp.Err = err.Error()
		jsonBytes, _ := json.Marshal(resp)
		w.Write(jsonBytes)
	} else {
		resp.Err = "unauthorized"
		jsonBytes, _ := json.Marshal(resp)
		w.Write(jsonBytes)
	}

}

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
		Username: user.Name,
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

//func GetUserIdFromToken(tokenReq string) (Id string) {
//	token, _ := jwt.Parse(tokenReq, func(token *jwt.Token) (interface{}, error) {
//		return []byte(""), nil
//	})
//	claim := token.Claims.(jwt.MapClaims)
//	id := claim["id"].(string)
//
//	return id
//}
