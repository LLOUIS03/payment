package auth

import "github.com/golang-jwt/jwt/v5"

var Secret = []byte("A7671BAC-8B17-44AC-A471-DDF2E0706F44")

type JwtCustomClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	ID    string `json:"merchant_id"`
}
