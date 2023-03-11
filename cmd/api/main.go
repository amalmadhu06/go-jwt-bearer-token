package main

import (
	//_ "github.com/amalmadhu06/go-jwt-bearer-token/cmd/api/docs"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/config"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/di"

	"log"
)

func main() {

	// @title Go Web API | JWT Bearer Token
	// @version 1.0
	// @description Go Web API to understand the implementation of jwt bearer token for authentication. Uses Gin and PostgreSQL. Follows clean architecture.

	// @contact.name API Support
	// @contact.url http://www.swagger.io/support
	// @contact.email support@swagger.io

	// @license.name MIT
	// @host localhost:3000
	// @license.url https://opensource.org/licenses/MIT

	// @BasePath /
	// @query.collection.format multi

	//Loading configs which are required for the working of our application
	cfg, configErr := config.LoadConfig()

	//In-case if there is an error loading config files
	if configErr != nil {
		log.Fatal("failed to load cgf", configErr)
	}

	//Calling InitializeAPI from di package.
	server, diErr := di.InitializeAPI(cfg)

	//In case if there is an error loading the dependencies
	if diErr != nil {
		log.Fatal("failed to start the server: ", diErr)
	} else {
		//Starting the server
		server.Start()
	}

}
