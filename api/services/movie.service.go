package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
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

		// Extract the substring between the last two slashes
		movURL := strings.TrimSuffix(movie.URL, "/")
		index := strings.LastIndex(movURL, "/")
		movURL = movURL[index+1:]
		// Convert the substring to an integer
		movieNumber, err := strconv.Atoi(movURL)
		if err != nil {
			fmt.Println("Error:", err)
			return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
		}

		var movieComments []models.Comment
		results := config.DB.Model(&models.Comment{}).Where("movie_id = ?", strconv.Itoa(movieNumber)).Find(&movieComments)
		if results.Error != nil {
			return (http.StatusBadGateway), gin.H{"status": "error", "message": results.Error}

		}

		movieResponses = append(movieResponses, responses.MovieResponse{
			Name:        movie.Title,
			OpeningCrawl: movie.OpeningCrawl,
			CommentCount: len(movieComments),
		})
	}

	cacheValue, err := json.Marshal(movieResponses)
	if err == nil {
		err := config.SetCache(cacheKey, string(cacheValue), 1*time.Minute)
		if err != nil {
			fmt.Println(err)
		}
	}

	return (http.StatusOK), gin.H{"status": "success", "data": movieResponses}
}
