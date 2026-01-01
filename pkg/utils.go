package pkg

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil

}

func ValidatePassword(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func ValidateEmail(email string) bool {
	return true
}

// JWTClaims defines the custom claims structure
type JWTClaims struct {
	UserID uint   `json:"userId"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string, email string) (string, error) {
	// Create JWT claims
	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
