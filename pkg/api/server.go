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

	// swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	user := engine.Group("/")
	{
		user.POST("/signup", userHandler.UserSignup)
		user.POST("/login", userHandler.UserLogin)
		user.POST("/logout", userHandler.UserLogout)
		user.POST("/access", userHandler.GetAccessToken)
		//use middleware here and write some routes which needs it
		user.Use(middleware.AuthorizeUser)
		{
			user.GET("/products", userHandler.ViewAllProducts)
		}

	}

	adminPanel := engine.Group("admin-panel")
	{
		adminPanel.POST("/login", adminHandler.AdminLogin)
		adminPanel.POST("/logout", adminHandler.AdminLogout)

		//use middleware here and write some routes which needs it
		adminPanel.Use(middleware.AuthorizeAdmin)
		{
			adminPanel.POST("/product", adminHandler.AddProduct)
			adminPanel.GET("/product", adminHandler.ViewAllProducts)
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
