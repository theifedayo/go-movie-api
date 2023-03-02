package responses

import (
	"github.com/theifedayo/go-movie-api/api/models"
)

// MovieResponse represents each movie data returned by the API
type MovieResponse struct {
	Title        string `json:"title"`
	OpeningCrawl string `json:"opening_crawl"`
	CommentCount int    `json:"comment_count"`
	//ReleaseDate  string `json:"release_date"`
}

// MovieResponse represents the movie data returned by the API as a List
type MovieListResponse struct {
	Results []models.Movie `json:"results"`
}
