package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/theifedayo/go-movie-api/api/controllers"
)

type MovieRouteController struct {
	movieController controllers.MovieController
}

func NewRouteMovieController(movieController controllers.MovieController) MovieRouteController {
	return MovieRouteController{movieController}
}

func (mc *MovieRouteController) MovieRoute(rg *gin.RouterGroup) {

	router := rg.Group("movies")

	router.GET("/", mc.movieController.ListMovies)
	router.POST("/:movieId/comments", mc.movieController.AddComment)
	router.GET("/:movieId/comments", mc.movieController.ListComments)

}
