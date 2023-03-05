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

// Movie Route routes all requests to /api/v1/movies
func (mc *MovieRouteController) MovieRoute(rg *gin.RouterGroup) {

	router := rg.Group("movies")
	//Route grouped to /api/v1/movies

	router.GET("/", mc.movieController.ListMovies)
	router.POST("/:movieId/comments", mc.movieController.AddCommentToMovies)
	router.GET("/:movieId/comments", mc.movieController.ListCommentsForAMovie)
	router.GET("/:movieId/characters", mc.movieController.GetCharactersForMovie)

}
