package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theifedayo/go-movie-api/models"
	"gorm.io/gorm"
)

type MovieController struct {
	DB *gorm.DB
}

func NewMovieController(DB *gorm.DB) MovieController {
	return MovieController{DB}
}

func (mc *MovieController) AddComment(ctx *gin.Context) {
	
}
