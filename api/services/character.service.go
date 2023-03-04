package services

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetCharactersForMovie(movieId string, ctx *gin.Context) (int, gin.H) {
	sortParam := ctx.Query("sort")
	filterParam := ctx.Query("gender")

	// Build the URL to fetch characters from swapi.dev API
	url := "https://swapi.dev/api/films/" + movieId + "/"

	// Add sorting and filtering parameters to the URL if specified
	if sortParam != "" {
		url += "?ordering=" + sortParam
	}
	if filterParam != "" {
		if strings.Contains(url, "?") {
			url += "&"
		} else {
			url += "?"
		}
		url += "gender=" + filterParam
	}

	// Fetch character data from swapi.dev API
	resp, err := http.Get(url)
	if err != nil {
		return (http.StatusInternalServerError), gin.H{"error": "Failed to fetch character data"}

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return (http.StatusNotFound), gin.H{"error": "Movie not found"}
	}

	// Parse the response and build the character list and metadata
	characters, totalHeightCm := parseCharacterData(resp.Body)
	totalHeightFeet, totalHeightInches := convertCmToFeetAndInches(float64(totalHeightCm))
	metadata := gin.H{
		"total_characters": len(characters),
		"total_height_cm":  totalHeightCm,
		"total_height_ft":  totalHeightFeet,
		"total_height_in":  totalHeightInches,
	}

	// Return the character list and metadata as the response
	return (http.StatusOK), gin.H{"characters": characters, "metadata": metadata}
}

func convertCmToFeetAndInches(cm float64) (float64, float64) {
	// 1 inch = 2.54 cm
	// 1 foot = 12 inches
	inches := cm / 2.54
	feet := math.Floor(inches / 12)
	inches = inches - (feet * 12)
	return feet, inches
}

// Helper function to parse the character data from the swapi.dev API response
func parseCharacterData(body io.Reader) ([]gin.H, int) {
	var data map[string]interface{}
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		return nil, 0
	}

	charactersData := data["characters"].([]interface{})
	var characters []gin.H
	var totalHeight int

	for _, characterData := range charactersData {
		characterURL := characterData.(string)
		characterResp, err := http.Get(characterURL)
		if err != nil {
			continue
		}

		characterBody, err := ioutil.ReadAll(characterResp.Body)
		if err != nil {
			continue
		}

		var character map[string]interface{}
		err = json.Unmarshal(characterBody, &character)
		if err != nil {
			continue
		}

		characters = append(characters, gin.H{
			"name":   character["name"],
			"gender": character["gender"],
			"height": character["height"],
		})

		// Add the character height to the total height
		characterHeight, err := strconv.Atoi(character["height"].(string))
		if err == nil {
			totalHeight += characterHeight
		}
	}

	return characters, totalHeight
}
