	// services/auth_service.go
	package services

	import (
		"errors"
		"os"
		"time"

		"github.com/golang-jwt/jwt/v5"
	)

	// Get JWT secret from environment variable, fallback to default (change in production!)
	var jwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "your-secret-key-change-this-in-production"))

	// Claims represents the JWT claims structure
	type Claims struct {
		UserID uint   `json:"userId"`
		Email  string `json:"email"`
		Role   string `json:"role"`
		jwt.RegisteredClaims
	}

	// GenerateToken generates a JWT token for a user
	func GenerateToken(userID uint, email, role string) (string, error) {
		if userID == 0 {
			return "", errors.New("invalid user ID")
		}
		if email == "" {
			return "", errors.New("email is required")
		}
		if role == "" {
			return "", errors.New("role is required")
		}

		// Create claims with user information
		claims := Claims{
			UserID: userID,
			Email:  email,
			Role:   role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    "visa-app",
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Sign token with secret key
		signedToken, err := token.SignedString(jwtSecret)
		if err != nil {
			return "", errors.New("failed to sign token")
		}

		return signedToken, nil
	}

	// ValidateToken validates a JWT token and returns the claims
	func ValidateToken(tokenString string) (*Claims, error) {
		if tokenString == "" {
			return nil, errors.New("token is required")
		}

		// Parse token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return jwtSecret, nil
		})

		if err != nil {
			return nil, errors.New("invalid or expired token")
		}

		// Extract claims
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}

		return nil, errors.New("invalid token claims")
	}

	// RefreshToken generates a new token for a user (extends expiration)
	func RefreshToken(oldToken string) (string, error) {
		// Validate old token
		claims, err := ValidateToken(oldToken)
		if err != nil {
			return "", err
		}

		// Generate new token with same user data
		return GenerateToken(claims.UserID, claims.Email, claims.Role)
	}

	// ExtractUserIDFromToken extracts the user ID from a token without full validation
	func ExtractUserIDFromToken(tokenString string) (uint, error) {
		claims, err := ValidateToken(tokenString)
		if err != nil {
			return 0, err
		}
		return claims.UserID, nil
	}

	// ExtractRoleFromToken extracts the role from a token without full validation
	func ExtractRoleFromToken(tokenString string) (string, error) {
		claims, err := ValidateToken(tokenString)
		if err != nil {
			return "", err
		}
		return claims.Role, nil
	}

	// IsTokenExpired checks if a token is expired
	func IsTokenExpired(tokenString string) bool {
		claims, err := ValidateToken(tokenString)
		if err != nil {
			return true
		}

		return claims.ExpiresAt.Time.Before(time.Now())
	}

	// GetTokenExpirationTime returns the expiration time of a token
	func GetTokenExpirationTime(tokenString string) (time.Time, error) {
		claims, err := ValidateToken(tokenString)
		if err != nil {
			return time.Time{}, err
		}

		return claims.ExpiresAt.Time, nil
	}

	// RevokeToken would typically mark a token as revoked in a database
	// For now, this is a placeholder that always returns an error
	// In production, you'd want to maintain a blacklist of revoked tokens
	func RevokeToken(tokenString string) error {
		// Validate token exists
		_, err := ValidateToken(tokenString)
		if err != nil {
			return err
		}

		// In a real implementation, you would:
		// 1. Store the token or its hash in a revoked tokens database table
		// 2. Check this table during token validation
		// 3. Clean up expired tokens periodically

		return errors.New("token revocation not implemented - use short expiration times instead")
	}

	// getEnvOrDefault gets an environment variable or returns a default value
	func getEnvOrDefault(key, defaultValue string) string {
		if value := os.Getenv(key); value != "" {
			return value
		}
		return defaultValue
	}
