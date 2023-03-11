package usecase

import (
	"fmt"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"
	services "github.com/amalmadhu06/go-jwt-bearer-token/pkg/usecase/interfaces"
	"github.com/golang-jwt/jwt/v4"
	"net/http"

	"time"
)

type JWTUseCase struct {
	SecretKey string
}

func NewJWTUserService() services.JWTUseCase {
	return &JWTUseCase{
		SecretKey: "secret",
		// Todo : replace this secretKey with your secret key. Ideally it should be fetched from .env file using viper
	}
}

func (j *JWTUseCase) GenerateRefreshToken(id int, email string, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":     id,
		"email":  email,
		"role":   role,
		"source": "refreshtoken",
		"exp":    time.Now().Local().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JWTUseCase) GenerateAccessToken(userid int, username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":     userid,
		"email":  username,
		"role":   role,
		"source": "accesstoken",
		"exp":    time.Now().Local().Add(time.Hour * 24).Unix(),
	}
	//todo : delete log after debugging
	fmt.Println("expiration time set in claims", claims["exp"])
	fmt.Printf("Data type of expiration : %T", claims["exp"])
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JWTUseCase) ParseToken(signedToken string) (*jwt.Token, *domain.JWTError) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &domain.JWTError{Code: http.StatusUnauthorized, Message: "unexpected signing method"}
		}
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		if validationErr, ok := err.(*jwt.ValidationError); ok {
			if validationErr.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, &domain.JWTError{
					Code:    http.StatusUnauthorized,
					Message: "token malformed",
				}
			} else if validationErr.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, &domain.JWTError{
					Code:    http.StatusUnauthorized,
					Message: "token expired or not yet valid",
				}
			} else {
				return nil, &domain.JWTError{
					Code:    http.StatusUnauthorized,
					Message: "token is invalid",
				}
			}
		} else {
			return nil, &domain.JWTError{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			}
		}
	}

	return token, nil
}

func (j *JWTUseCase) VerifyToken(signedToken string) (bool, *domain.SignedDetails, error) {
	token, err := j.ParseToken(signedToken)
	if err != nil {
		return false, nil, err
	}

	if !token.Valid {
		return false, nil, &domain.JWTError{
			Code:    http.StatusUnauthorized,
			Message: "token is invalid",
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil, &domain.JWTError{
			Code:    http.StatusInternalServerError,
			Message: "error parsing token claims",
		}
	}

	signedDetails := &domain.SignedDetails{
		ID:     int(claims["id"].(float64)),
		Email:  claims["email"].(string),
		Role:   claims["role"].(string),
		Source: claims["source"].(string),
		Exp:    int64(claims["exp"].(float64)),
		//ExpiresAt: time.Unix(int64(claims["exp"].(float64)), 0),
	}

	if err := signedDetails.Valid(); err != nil {
		return false, nil, err
	}

	return true, signedDetails, nil
}
