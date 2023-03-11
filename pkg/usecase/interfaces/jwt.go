package interfaces

import (
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
	"github.com/golang-jwt/jwt/v4"
)

type JWTUseCase interface {
	//GenerateRefreshToken(userid int, username string, role string) (string, error)
	//GenerateAccessToken(userid int, username string, role string) (string, error)
	//VerifyToken(signedToken string) (bool, *domain.SignedDetails)
	//GetTokenFromString(signedToken string, claims *domain.SignedDetails) (*jwt.Token, error)

	GenerateRefreshToken(id int, email string, role string) (string, error)
	GenerateAccessToken(userid int, username string, role string) (string, error)
	ParseToken(signedToken string) (*jwt.Token, *domain.JWTError)
	VerifyToken(signedToken string) (bool, *domain.SignedDetails, error)
}
