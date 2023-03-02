package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theifedayo/go-movie-api/api/models"
	"github.com/theifedayo/go-movie-api/config"
	"gorm.io/gorm"
)

type MovieController struct {
	DB *gorm.DB
}

func NewMovieController(DB *gorm.DB) MovieController {
	return MovieController{DB}
}

func (mc *MovieController) AddComment(ctx *gin.Context) {
	movieId := ctx.Param("movieId")
	var payload *models.CreateCommentRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()
	newComment := models.Comment{
		MovieID:   movieId,
		Comment:   payload.Comment,
		IP:        config.GetIPAddress(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := config.DB.Create(&newComment)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newComment})
}

func (mc *MovieController) ListComments(ctx *gin.Context) {
	movieId := ctx.Param("movieId")
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var comments []models.Comment
	results := config.DB.Order("created_at desc").Limit(intLimit).Offset(offset).Where("movie_id = ?", movieId).Find(&comments)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(comments), "data": comments})
}
