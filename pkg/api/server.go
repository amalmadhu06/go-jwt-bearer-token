package http

import (
	"fmt"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/api/handler"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/api/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(userHandler *handler.UserHandler, adminHandler *handler.AdminHandler, middleware middleware.Middleware) *ServerHTTP {
	engine := gin.New()

	//	logger logs following info for each request : http method, req URL, remote address of the client, res status code, elapsed time
	engine.Use(gin.Logger())

	//	handler for generating swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	user := engine.Group("/")
	{
		user.POST("/signup", userHandler.UserSignup)
		user.POST("/login", userHandler.UserLogin)
		user.GET("/logout", userHandler.UserLogout)

		//use middleware here and write some routes which needs it
		user.Use(middleware.AuthorizeJWT)
		{
			user.GET("/view-all-products", userHandler.ViewAllProducts)
		}

	}

	adminPanel := engine.Group("adminPanel")
	{
		adminPanel.POST("/login", adminHandler.AdminLogin)
		adminPanel.GET("/logout", adminHandler.AdminLogout)

		//use middleware here and write some routes which needs it
		adminPanel.Use(middleware.AuthorizeJWT)
		{
			adminPanel.POST("/add-new-product", adminHandler.AddProduct)
			adminPanel.GET("/view-all-product", adminHandler.ViewAllProducts)
		}

	}

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	err := sh.engine.Run(":3000")
	if err != nil {
		fmt.Println("error starting server :", err.Error())
		return
	}
}
