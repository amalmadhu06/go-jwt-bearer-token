package domain

import "github.com/golang-jwt/jwt/v4"

// SignedDetails represents the JWT token details
type SignedDetails struct {
	ID     int    `json:"id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Source string `json:"source"`
	Exp    int64  `json:"exp"`

	// Custom claims as per application needs
	jwt.StandardClaims
}

// JWTError represents an error related to JWT token processing
type JWTError struct {
	Code    int
	Message string
}

func (e *JWTError) Error() string {
	return e.Message
}
