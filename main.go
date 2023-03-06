package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/theifedayo/go-movie-api/api/controllers"
	"github.com/theifedayo/go-movie-api/api/routes"
	"github.com/theifedayo/go-movie-api/config"
	_ "github.com/theifedayo/go-movie-api/docs"
)

var (
	server               *gin.Engine
	MovieController      controllers.MovieController
	MovieRouteController routes.MovieRouteController
)

func init() {
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("could not load environment variables", err)
	}
	//Initialize the database connection
	config.ConnectToDB(&configs)
	//Initialize the Redis client
	config.SetRedisConfig(&configs)

	server = gin.Default()
}

// @title Go Movie API
// @version 1.0
// @description This is a RESTful API that provides information about Star Wars movies.
// @host http://gomovie-api.herokuapp.com
// @BasePath /api/v1
func main() {
	configs, _ := config.LoadConfig(".")

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", configs.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))
	//redirect base url to swagger documentation
	server.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(301, "/api/v1/docs/index.html")
	})

	router := server.Group("/api/v1")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Welcome to Movie API"})
	})

	//register the Swagger route and Swagger UI route
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	MovieRouteController.MovieRoute(router)

	log.Fatal(server.Run(":" + configs.ServerPort))
}
