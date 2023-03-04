package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/theifedayo/go-movie-api/api/models"
	"github.com/theifedayo/go-movie-api/api/responses"
	"github.com/theifedayo/go-movie-api/config"
)

func ListMovies(ctx *gin.Context) (int, gin.H) {

	var movies responses.MovieListResponse
	var movieResponses []responses.MovieResponse

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
		return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
	}

	resp, err := http.Get("https://swapi.dev/api/films/")
	if err != nil {
		return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
	}

	fmt.Println(json.NewDecoder(resp.Body))

	err = json.NewDecoder(resp.Body).Decode(&movies)
	if err != nil {
		return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
	}

	sort.Slice(movies.Results, func(i, j int) bool {
		iDate, _ := time.Parse("2006-01-02", movies.Results[i].ReleaseDate)
		jDate, _ := time.Parse("2006-01-02", movies.Results[j].ReleaseDate)
		return iDate.Before(jDate)
	})

	for _, movie := range movies.Results {

		// Extract the substring between the last
		// two slashes from the movie url to get movieId
		movieURL := strings.TrimSuffix(movie.URL, "/")
		index := strings.LastIndex(movieURL, "/")
		movieNumber := movieURL[index+1:]

		if err != nil {
			fmt.Println("Error:", err)
			return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
		}

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

func GetMovie(movieID string) (*models.Movie, error) {
	var movie models.Movie

	// First, try to fetch the movie from the database
	err := config.DB.Where("id = ?", movieID).First(&movie).Error
	if err == nil {
		return &movie, nil
	}

	// If the movie is not found in the database, fetch it from the external API
	res, err := http.Get(fmt.Sprintf("https://swapi.dev/api/films/%s/", movieID))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch movie from external API: %s", res.Status)
	}
	return &movie, nil
}
