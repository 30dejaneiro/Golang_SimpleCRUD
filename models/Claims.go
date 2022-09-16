package models

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTService interface {
	GenerateToken(username string, isAdmin bool) string
	ValidateToken(tokenString string) (*jwt.Token, error)
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: "JWT_SECRET",
	}
}

type JwtClaims struct {
	Username string `json:"username,omitempty"`
	IsAdmin  bool   `json:"isAdmin,omitempty"`
	jwt.StandardClaims
}
type jwtService struct {
	secretKey string
}

func (jwtSrv *jwtService) GenerateToken(username string, isAdmin bool) string {
	// Set custom and standard claims

	claims := &JwtClaims{
		username,
		isAdmin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 10).Unix(),
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token using the secret signing key
	t, err := token.SignedString([]byte(jwtSrv.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (jwtSrv *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Signing method validation
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret signing key
		return []byte(jwtSrv.secretKey), nil
	})
}
