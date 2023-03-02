package controllers

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
	"gorm.io/gorm"
)

type MovieController struct {
	DB *gorm.DB
}

func NewMovieController(DB *gorm.DB) MovieController {
	return MovieController{DB}
}

func (mc *MovieController) ListMovies(ctx *gin.Context) {

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
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": movieResponses})
			return
		}

		fmt.Println(err)
	}

	resp, err := http.Get("https://swapi.dev/api/films/")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fmt.Println(json.NewDecoder(resp.Body))

	err = json.NewDecoder(resp.Body).Decode(&movies)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
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
			return
		}

		var movieComments []models.Comment
		results := config.DB.Model(&models.Comment{}).Where("movie_id = ?", strconv.Itoa(movieNumber)).Find(&movieComments)
		if results.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
			return
		}

		movieResponses = append(movieResponses, responses.MovieResponse{
			Title:        movie.Title,
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

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": movieResponses})

}
