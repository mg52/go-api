package domain

import "github.com/golang-jwt/jwt/v4"

type Authentication struct {
	Token string `json:"token"`
}

type Login struct {
	Username string `json:"username" validate:"required,min=3,max=10"`
	Password string `json:"password" validate:"required,min=3,max=12"`
}

type Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
