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

// @Title Go Movie API
func main() {
	configs, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", configs.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api/v1")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Movie API"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	//register the Swagger route and Swagger UI route
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	MovieRouteController.MovieRoute(router)

	log.Fatal(server.Run(":" + configs.ServerPort))
}
