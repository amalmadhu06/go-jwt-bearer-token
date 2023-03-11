//package middleware
//
//import (
//	"errors"
//	"fmt"
//	"net/http"
//	"strings"
//	"time"
//
//	"github.com/gin-gonic/gin"
//
//	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/usecase/interfaces"
//)
//
//const (
//	authHeaderKey   = "Authorization"
//	accessTokenType = "accesstoken"
//)
//
//type Middleware interface {
//	AuthorizeJWT(*gin.Context)
//}
//
//type middlewareImpl struct {
//	jwtService interfaces.JWTUseCase
//}
//
////func NewUserMiddleware(jwtService interfaces.JWTUseCase) Middleware {
////	return &middlewareImpl{
////		jwtService: jwtService,
////	}
////}
//
//func NewAdminMiddleware(jwtService interfaces.JWTUseCase) Middleware {
//	return &middlewareImpl{
//		jwtService: jwtService,
//	}
//}
//
//func (mw *middlewareImpl) AuthorizeJWT(c *gin.Context) {
//	// Get the Authorization header value
//	authHeader := c.GetHeader(authHeaderKey)
//	if authHeader == "" {
//		err := errors.New("request does not contain an access token")
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//	//todo : delete after debugging
//	fmt.Println("auth header", authHeader)
//
//	// Extract the access token from the header
//	bearerToken := strings.Split(authHeader, " ")
//	if len(bearerToken) != 2 {
//		err := errors.New("invalid access token format")
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//
//	accessToken := bearerToken[1]
//	//todo : delete after debugging
//	fmt.Println("access token", accessToken)
//	// Verify the access token
//	ok, claims, err := mw.jwtService.VerifyToken(accessToken)
//	//todo : delete after debugging
//	fmt.Println("ok: ", ok, "claims : ", claims, "err :", err)
//	if err != nil {
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//
//	if !ok {
//		err := errors.New("invalid access token")
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//
//	// Check if the token type is valid
//	if claims.Source != accessTokenType {
//		err := errors.New("invalid token type")
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//
//	// Check if the token is expired
//	//todo : delete this after debugging
//	fmt.Println(claims["exp"])
//
//	if time.Now().Unix() > claims["exp"] {
//		err := errors.New("access token has expired")
//		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
//			"error": err.Error(),
//		})
//		return
//	}
//
//	// Set the user email and ID in the context for future use
//	c.Set("user_email", claims.Email)
//	c.Set("user_id", claims.Id)
//
//	c.Next()
//}

package middleware

import (
	"errors"
	"fmt"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/usecase/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const (
	authHeaderKey   = "Authorization"
	accessTokenType = "accesstoken"
)

type Middleware interface {
	AuthorizeJWT(*gin.Context)
}

type middlewareImpl struct {
	jwtService interfaces.JWTUseCase
}

func NewAdminMiddleware(jwtService interfaces.JWTUseCase) Middleware {
	return &middlewareImpl{
		jwtService: jwtService,
	}
}

func (mw *middlewareImpl) AuthorizeJWT(c *gin.Context) {
	// Get the Authorization header value
	authHeader := c.GetHeader(authHeaderKey)
	if authHeader == "" {
		err := errors.New("request does not contain an access token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Extract the access token from the header
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		err := errors.New("invalid access token format")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	accessToken := bearerToken[1]

	// Verify the access token
	ok, claims, err := mw.jwtService.VerifyToken(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !ok {
		err := errors.New("invalid access token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if the token type is valid
	if claims.Source != accessTokenType {
		err := errors.New("invalid token type")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if the token is expired
	//todo : delete after debugging
	fmt.Println(claims.Exp)
	if time.Now().Unix() > claims.Exp {
		err := errors.New("access token has expired")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Set the user email and ID in the context for future use
	c.Set("user_email", claims.Email)
	c.Set("user_id", claims.Id)
	c.Set("user_role", claims.Role)

	c.Next()
}
