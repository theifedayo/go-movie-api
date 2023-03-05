package responses

import (
	"github.com/theifedayo/go-movie-api/api/models"
)

// MovieResponse represents each movie data returned by the API
type MovieResponse struct {
	Name         string `json:"name"`
	OpeningCrawl string `json:"opening_crawl"`
	CommentCount int    `json:"comment_count"`
	//ReleaseDate  string `json:"release_date"`
}

// MovieResponse represents the movie data returned by the API as a List
type MovieListResponse struct {
	Results []models.Movie `json:"results"`
}

// CharacterMetadata represents the meta data returned by the API for total characters
type CharacterMetadata struct {
	TotalCharacters int     `json:"total_characters"`
	TotalHeightCm   float64 `json:"total_height_cm"`
	TotalHeightFt   float64 `json:"total_height_ft"`
	TotalHeightIn   float64 `json:"total_height_"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AddCommentRequest struct {
	Comment string `json:"comment" binding:"required"`
}

type AddCommentResponse struct {
	Data models.CreateCommentRequest `json:"data"`
}

// CharacterResponse represents the response for getting a list of characters for a movie
type CharacterResponse struct {
	Characters []string `json:"characters"`
}
