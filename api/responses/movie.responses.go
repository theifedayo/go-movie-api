package responses

import (
	"github.com/theifedayo/go-movie-api/api/models"
)

// MovieResponse represents the movie data returned by the API
type MovieResponse struct {
	Title        string `json:"title"`
	OpeningCrawl string `json:"opening_crawl"`
	CommentCount int    `json:"comment_count"`
	ReleaseDate  string `json:"release_date"`
}

type MovieListResponse struct {
	Results []models.Movie `json:"results"`
}

// CommentResponse represents the comment data returned by the API
type CommentResponse struct {
}
