package middleware

import (
	"errors"
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
	AuthorizeAdmin(c *gin.Context)
	AuthorizeUser(c *gin.Context)
}

type middlewareImpl struct {
	jwtService interfaces.JWTUseCase
}

func NewAdminMiddleware(jwtService interfaces.JWTUseCase) Middleware {
	return &middlewareImpl{
		jwtService: jwtService,
	}
}

func isTokenValid(c *gin.Context, mw *middlewareImpl) bool {
	// Get the Authorization header value
	authHeader := c.GetHeader(authHeaderKey)
	if authHeader == "" {
		err := errors.New("request does not contain an access token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return false
	}

	// Extract the access token from the header
	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 {
		err := errors.New("invalid access token format")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return false
	}

	accessToken := bearerToken[1]

	// Verify the access token
	ok, claims, err := mw.jwtService.VerifyToken(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return false
	}

	if !ok {
		err := errors.New("invalid access token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return false
	}
	// Check if the token type is valid
	if claims.Source != accessTokenType {
		err := errors.New("invalid token type")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return false
	}

	// Check if the token is expired
	if time.Now().Unix() > claims.Exp {
		err := errors.New("access token has expired")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return false
	}

	// Set the user email and ID in the context for future use
	c.Set("user_email", claims.Email)
	c.Set("user_id", claims.Id)
	c.Set("user_role", claims.Role)
	return true
}

func (mw *middlewareImpl) AuthorizeUser(c *gin.Context) {
	if !isTokenValid(c, mw) {
		return
	}

	// Check if the user role is valid
	role, _ := c.Get("user_role")
	if role != "user" {
		err := errors.New("user is not authorized")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Next()
}

func (mw *middlewareImpl) AuthorizeAdmin(c *gin.Context) {
	if !isTokenValid(c, mw) {
		return
	}

	// Check if the user role is valid
	role, _ := c.Get("user_role")
	if role != "admin" {
		err := errors.New("admin is not authorized")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Next()
}

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
//
//	// Verify the access token
//	ok, claims, err := mw.jwtService.VerifyToken(accessToken)
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
//	if time.Now().Unix() > claims.Exp {
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
//	c.Set("user_role", claims.Role)
//	c.Next()
//}
