package services

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theifedayo/go-movie-api/api/models"
	"github.com/theifedayo/go-movie-api/config"
)


// AddCommentToMovies adds a new comment for a movie.
// It takes a Context of ctx and the Id of the movie to add comment as a parameter and returns a status code, as well as a map containing necessary information.
// It also returns error status code and a map of error message if one occurs
func AddCommentToMovies(movieId string, ctx *gin.Context) (int, gin.H) {
	var payload *models.CreateCommentRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {

		return http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()}
	}

	now := time.Now().UTC().Truncate(time.Second)
	newComment := models.Comment{
		MovieID:   movieId,
		Comment:   payload.Comment,
		IP:        config.GetIPAddress(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := config.DB.Create(&newComment)
	if result.Error != nil {
		return http.StatusBadRequest, gin.H{"status": "error", "message": result.Error.Error()}
	}
	return (http.StatusCreated), gin.H{"status": "success", "data": newComment}
}


// ListCommentsForAMovie lists the comments for a specific movie.
// It takes a Context of ctx and the Id of the movie to get all its comments as a parameter and returns a status code, as well as a map containing necessary information.
// It also returns error status code and a map of error message if one occurs
func ListCommentsForAMovie(movieId string, ctx *gin.Context) (int, gin.H) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var comments []models.Comment
	results := config.DB.Order("created_at desc").Limit(intLimit).Offset(offset).Where("movie_id = ?", movieId).Find(&comments)
	if results.Error != nil {
		return (http.StatusBadGateway), gin.H{"status": "error", "message": results.Error}

	}

	return (http.StatusOK), gin.H{"status": "success", "results": len(comments), "data": comments}
}
