package domain

import "github.com/golang-jwt/jwt/v4"

type Authentication struct {
	Token string `json:"token"`
}

type Login struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
