package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/theifedayo/go-movie-api/api/services"
	"gorm.io/gorm"
)

type MovieController struct {
	DB *gorm.DB
}

func NewMovieController(DB *gorm.DB) MovieController {
	return MovieController{DB}
}

// @Summary Get a list of movies
// @Description Retrieves a list of movies sorted by release date, along with opening crawls and comment count
// @Accept  json
// @Produce  json
// @Success 200 {object} MoviesResponse
// @Router /movies [get]
func (mc *MovieController) ListMovies(ctx *gin.Context) {
	statusCode, result := services.ListMovies(ctx)
	ctx.JSON(statusCode, result)
}

func (mc *MovieController) AddCommentToMovies(ctx *gin.Context) {
	movieId := ctx.Param("movieId")
	statusCode, result := services.AddCommentToMovies(movieId, ctx)
	ctx.JSON(statusCode, result)
}

func (mc *MovieController) ListCommentsForAMovie(ctx *gin.Context) {
	movieId := ctx.Param("movieId")
	statusCode, result := services.ListCommentsForAMovie(movieId, ctx)
	ctx.JSON(statusCode, result)
}

func (mc *MovieController) GetCharactersForMovie(ctx *gin.Context) {
	movieId := ctx.Param("movieId")
	gender := ctx.Query("gender")
	sort := ctx.Query("sort")
	order := ctx.Query("order")
	statusCode, result := services.GetCharactersForMovie(movieId, sort, order, gender, ctx)
	ctx.JSON(statusCode, result)
}
