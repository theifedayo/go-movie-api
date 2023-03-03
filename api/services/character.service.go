package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/theifedayo/go-movie-api/api/models"
)

func GetCharactersForMovie(movieID string, ctx *gin.Context) (int, gin.H) {
	movieId := ctx.Param("movieId")
	sortBy := ctx.DefaultQuery("sort_by", "name")
	order := ctx.DefaultQuery("order", "asc")
	genderFilter := ctx.DefaultQuery("gender", "")

	// Fetch the movie from the database or the external API
	//movie, err := models.GetMovie(movieId)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// var movieModel models.Movie
	// movie := config.DB.Order("created_at desc").Where("movie_id = ?", movieId).Find(&movieModel)

	// Fetch the characters from the external API
	var characters []models.Character
	url := fmt.Sprintf("https://swapi.dev/api/films/%s/", movieId)
	for {
		resp, err := http.Get(url)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var data struct {
			Characters []string `json:"characters"`
			Next       string   `json:"next"`
		}
		if err := json.Unmarshal(body, &data); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, characterURL := range data.Characters {
			character, err := models.GetCharacter(characterURL)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Filter characters by gender if specified
			if genderFilter != "" && character.Gender != genderFilter {
				continue
			}

			characters = append(characters, character)
		}

		if data.Next == "" {
			break
		}
		url = data.Next
	}

	// Sort characters by the specified field and order
	switch strings.ToLower(sortBy) {
	case "name":
		if strings.ToLower(order) == "asc" {
			sort.Slice(characters, func(i, j int) bool {
				return characters[i].Name < characters[j].Name
			})
		} else {
			sort.Slice(characters, func(i, j int) bool {
				return characters[i].Name > characters[j].Name
			})
		}
	case "gender":
		if strings.ToLower(order) == "asc" {
			sort.Slice(characters, func(i, j int) bool {
				return characters[i].Gender < characters[j].Gender
			})
		} else {
			sort.Slice(characters, func(i, j int) bool {
				return characters[i].Gender > characters[j].Gender
			})
		}
	case "height":
		if strings.ToLower(order) == "asc" {
			sort.Slice(characters, func(i, j int) bool {
				return characters[i].Height < characters[j].Height
			})
		} else {
			sort.Slice(characters, func(i, j int) bool {
				return characters[i].Height > characters[j].Height
			})
		}
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort_by parameter"})
		return
	}

	totalHeightCm := 0
	for _, character := range characters {
		height, err := strconv.Atoi(character.Height)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		totalHeightCm += height
	}
	totalHeightInch := float64(totalHeightCm) / 2.54
	totalHeightFeet := int(totalHeightInch / 12)
	totalHeightInch = totalHeightInch - float64(totalHeightFeet*12)
	totalHeight := fmt.Sprintf("%dft %.2fin (%dcm)", totalHeightFeet, totalHeightInch, totalHeightCm)

	// Return the response
	ctx.JSON(http.StatusOK, gin.H{
		"metadata": gin.H{
			"total_count": len(characters),
			"total_height": gin.H{
				"cm":   totalHeightCm,
				"feet": totalHeightFeet,
				"inch": totalHeightInch,
				"desc": totalHeight,
			},
		},
		"data": characters,
	})
	return
}