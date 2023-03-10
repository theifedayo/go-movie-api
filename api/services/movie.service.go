package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theifedayo/go-movie-api/api/models"
	"github.com/theifedayo/go-movie-api/api/responses"
	"github.com/theifedayo/go-movie-api/config"
)

// ListMovies checks for cached data, if not fetches movie data online from https://swapi.dev/api/films, sorts them by release date, caches and lists them .
// It takes a Context of ctx as a parameter and returns a status code, as well as a map containing necessary information.
// It also returns error status code and a map of error message if one occurs
func ListMovies(ctx *gin.Context) (int, gin.H) {

	var movies responses.MovieListResponse
	var movieResponses []responses.MovieResponse

	//check if the result has been cached
	cacheKey := "movies"
	cacheResult, err := config.GetCache(cacheKey, 5*time.Minute)
	if err != nil {
		fmt.Println(err)
	}

	if cacheResult != "" {
		err := json.Unmarshal([]byte(cacheResult), &movieResponses)
		if err == nil {
			return (http.StatusOK), gin.H{"status": "success", "data": movieResponses}

		}

		fmt.Println(err)
		return (http.StatusInternalServerError), gin.H{"status": "error", "message": err.Error()}
	}

	//if not, make a GET request
	resp, err := http.Get("https://swapi.dev/api/films/")
	if err != nil {
		return (http.StatusInternalServerError), gin.H{"status": "error", "message": err.Error()}

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return (http.StatusInternalServerError), gin.H{"status": "error", "message": err.Error()}
	}

	fmt.Println(json.NewDecoder(resp.Body))

	err = json.NewDecoder(resp.Body).Decode(&movies)
	if err != nil {
		return (http.StatusInternalServerError), gin.H{"status": "error", "message": err.Error()}
	}

	for _, movie := range movies.Results {

		// Extract the substring between the last two slashes from the movie url to get movieId
		// and get all comments for the movie
		movieURL := strings.TrimSuffix(movie.URL, "/")
		index := strings.LastIndex(movieURL, "/")
		movieNumber := movieURL[index+1:]

		var movieComments []models.Comment
		results := config.DB.Model(&models.Comment{}).Where("movie_id = ?", movieNumber).Find(&movieComments)
		if results.Error != nil {
			return (http.StatusBadGateway), gin.H{"status": "error", "message": results.Error}

		}

		movieResponses = append(movieResponses, responses.MovieResponse{
			Name:         movie.Title,
			OpeningCrawl: movie.OpeningCrawl,
			CommentCount: len(movieComments),
		})
	}

	//cache the movie response
	cacheValue, err := json.Marshal(movieResponses)
	if err == nil {
		//cache expiration time set to 5 minutes
		err := config.SetCache(cacheKey, string(cacheValue), 5*time.Minute)
		if err != nil {
			fmt.Println(err)
		}
	}

	return (http.StatusOK), gin.H{"status": "success", "data": movieResponses}
}
