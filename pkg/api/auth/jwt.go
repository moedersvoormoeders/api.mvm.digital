package auth

import (
	"github.com/dgrijalva/jwt-go"
)

type Claim struct {
	Name string `json:"name"`
	// TODO: roles
	jwt.StandardClaims
}
