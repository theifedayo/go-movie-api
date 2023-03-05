package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/theifedayo/go-movie-api/api/responses"
	"github.com/theifedayo/go-movie-api/api/services"
	"gorm.io/gorm"
)

type MovieController struct {
	DB *gorm.DB
}

func NewMovieController(DB *gorm.DB) MovieController {
	return MovieController{DB}
}

var MovieResponse responses.MovieResponse

// @Summary Get a list of movies
// @Description Retrieves a list of movies sorted by release date, along with name, opening crawls and comment count
// @Tags Movie
// @Produce  json
// @Success 200 {object} responses.MovieResponse
// @Failure 500  {object} responses.ErrorResponse
// @Router /movies [get]
func (mc *MovieController) ListMovies(ctx *gin.Context) {
	statusCode, result := services.ListMovies(ctx)
	ctx.JSON(statusCode, result)
}

// @Summary Add a new comment
// @Description Add a new comment for the specified movie
// @Tags Movie
// @Accept json
// @Produce json
// @Param movieId path int true "Movie ID"
// @Param comment body responses.AddCommentRequest true "Comment request body"
// @Success 201 {object} responses.AddCommentResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /movies/{movieId}/comments [post]
func (mc *MovieController) AddCommentToMovies(ctx *gin.Context) {
	movieId := ctx.Param("movieId")
	statusCode, result := services.AddCommentToMovies(movieId, ctx)
	ctx.JSON(statusCode, result)
}

// @Summary Get a list of movies
// @Description Returns a list of comments for the specified movie
// @Tags Movie
// @Produce  json
// @Param movieId path int true "Movie ID"
// @Failure 404  {object} responses.ErrorResponse
// @Failure 500  {object} responses.ErrorResponse
// @Router /movies/{movieId}/comments [get]
func (mc *MovieController) ListCommentsForAMovie(ctx *gin.Context) {
	movieId := ctx.Param("movieId")
	statusCode, result := services.ListCommentsForAMovie(movieId, ctx)
	ctx.JSON(statusCode, result)
}

// @Summary Get characters for a movie
// @Description Returns a list of characters for the specified movie
// @Tags Movie
// @Param movieId path int true "Movie ID"
// @Param sort query string false "[Optional] The field to sort the characters by one of name, gender, or height"
// @Param order query string false "[Optional] Use asc or desc to sort in ascending or descending order, respectively. For example, ?sort=height&order=desc will sort by height in descending order, while ?sort=height&order=asc will sort by height in ascending order"
// @Param gender query string false "[Optional] The filter criteria to apply to the characters to filter by male or female. For example, ?gender=male will filter by male characters and return only male characters and ?sort=height&order=desc&gender=female will filter by female characters, listing only female characters with their height in descending order"
// @Success 200 {object} responses.CharacterResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /movies/{movieId}/characters [get]
func (mc *MovieController) GetCharactersForMovie(ctx *gin.Context) {
	movieId := ctx.Param("movieId")
	gender := ctx.Query("gender")
	sort := ctx.Query("sort")
	order := ctx.Query("order")
	statusCode, result := services.GetCharactersForMovie(movieId, sort, order, gender, ctx)
	ctx.JSON(statusCode, result)
}
