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
	"github.com/theifedayo/go-movie-api/api/responses"
)

// GetCharactersForMovie gets the list of characters a specific movie.
// It takes a Context of ctx and the Id of the movie to get all its comments as a parameter and returns a status code, as well as a map containing necessary information.
// It can also accept sort parameters to sort by one of name, gender or height in ascending or descending order and filter parameter to filter by gender.
// It also returns error status code and a map of error message if one occurs
func GetCharactersForMovie(movieId string, sortParam string, order string, gender string, ctx *gin.Context) (int, gin.H) {
	var characters []models.Character
	var metadata responses.CharacterMetadata

	url := fmt.Sprintf("https://swapi.dev/api/films/%s/", movieId)

	res, err := http.Get(url)
	if err != nil {
		return (http.StatusInternalServerError), gin.H{"status": "error", "message": err.Error()}
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return (http.StatusInternalServerError), gin.H{"status": "error", "message": err.Error()}
	}

	var movieData map[string]interface{}
	err = json.Unmarshal(body, &movieData)
	if err != nil {
		return (http.StatusInternalServerError), gin.H{"status": "error", "message": err.Error()}
	}

	charactersUrls := movieData["characters"].([]interface{})

	metadata.TotalCharacters = len(charactersUrls)

	var totalHeightCm float64 = 0

	for _, characterUrl := range charactersUrls {
		res, err := http.Get(characterUrl.(string))
		if err != nil {
			return (http.StatusInternalServerError), gin.H{"status": "error", "message": err.Error()}
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return (http.StatusInternalServerError), gin.H{"status": "error", "message": err.Error()}
		}

		var characterData map[string]interface{}
		err = json.Unmarshal(body, &characterData)
		if err != nil {
			return (http.StatusInternalServerError), gin.H{"status": "error", "message": err.Error()}
		}

		heightCm, err := strconv.ParseFloat(strings.Replace(characterData["height"].(string), ",", "", -1), 64)
		if err != nil {
			return (http.StatusInternalServerError), gin.H{"status": "error", "message": err.Error()}
		}

		totalHeightCm += heightCm

		characters = append(characters, models.Character{
			Name:   characterData["name"].(string),
			Height: characterData["height"].(string),
			Gender: characterData["gender"].(string),
		})
	}

	// Filter characters by gender, if gender filter is provided
	if gender != "" {
		totalHeightCm = 0
		characters, totalHeightCm = sortByGender(characters, gender)
	} else {
		var filteredCharacters []models.Character
		filteredCharacters = characters
		fmt.Println(filteredCharacters)
	}

	// Sort characters based on given sort parameter
	switch sortParam {
	case "name":
		if order == "desc" {
			characters = sortDescByName(characters)
		} else {
			characters = sortAscByName(characters)
		}
	case "gender":
		if order == "desc" {
			characters = sortDescByGender(characters)
		} else {
			characters = sortAscByGender(characters)
		}
	case "height":
		if order == "desc" {
			characters = sortDescByHeight(characters)
		} else {
			characters = sortAscByHeight(characters)
		}
	}

	metadata.TotalHeightCm = totalHeightCm
	metadata.TotalHeightFt = cmToFeet(totalHeightCm)
	metadata.TotalHeightIn = cmToInch(totalHeightCm)

	return (http.StatusOK), gin.H{
		"status": "success",
		"metadata": gin.H{
			"total_count": len(characters),
			"total_height": gin.H{
				"cm":   metadata.TotalHeightCm,
				"feet": metadata.TotalHeightFt,
				"inch": metadata.TotalHeightIn,
			},
		},
		"data": characters,
	}

}

//** Character Service Helpers**//

func sortByGender(characters []models.Character, gender string) ([]models.Character, float64) {
	var filteredCharacters []models.Character
	var totalHeightCm float64 = 0
	for _, character := range characters {
		if strings.ToLower(character.Gender) == strings.ToLower(gender) {
			filteredCharacters = append(filteredCharacters, character)
		}
	}

	for _, filteredCharacters := range filteredCharacters {
		heightCm, err := strconv.ParseFloat(strings.Replace(filteredCharacters.Height, ",", "", -1), 64)
		if err != nil {
			return characters, heightCm
		}
		totalHeightCm += heightCm
	}
	sort.Slice(filteredCharacters, func(i, j int) bool {
		return filteredCharacters[i].Name < filteredCharacters[j].Name
	})
	characters = filteredCharacters
	return characters, totalHeightCm
}

func sortDescByName(characters []models.Character) []models.Character {
	sort.Slice(characters, func(i, j int) bool {
		return characters[i].Name > characters[j].Name
	})
	return characters
}

func sortAscByName(characters []models.Character) []models.Character {
	sort.Slice(characters, func(i, j int) bool {
		return characters[i].Name < characters[j].Name
	})
	return characters
}

func sortDescByGender(characters []models.Character) []models.Character {
	sort.Slice(characters, func(i, j int) bool {
		return characters[i].Gender > characters[j].Gender
	})
	return characters
}

func sortAscByGender(characters []models.Character) []models.Character {
	sort.Slice(characters, func(i, j int) bool {
		return characters[i].Gender < characters[j].Gender
	})
	return characters
}

func sortDescByHeight(characters []models.Character) []models.Character {
	sort.Slice(characters, func(i, j int) bool {
		return characters[i].Height > characters[j].Height
	})
	return characters
}

func sortAscByHeight(characters []models.Character) []models.Character {
	sort.Slice(characters, func(i, j int) bool {
		return characters[i].Height < characters[j].Height
	})
	return characters
}

func cmToFeet(cm float64) float64 {
	feet := cm / 30.48
	return RoundToTwoDecimalPlaces(feet)
}

func cmToInch(cm float64) float64 {
	inches := cm / 2.54
	return RoundToTwoDecimalPlaces(inches)
}

func RoundToTwoDecimalPlaces(num float64) float64 {
	return float64(int(num*100)) / 100
}
