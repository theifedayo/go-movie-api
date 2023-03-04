package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/theifedayo/go-movie-api/api/models"
	"github.com/theifedayo/go-movie-api/api/responses"
)

func GetCharactersForMovie(movieId string, sortParam string, order string, gender string, ctx *gin.Context) (int, gin.H) {
	var characters []models.Character
	var metadata responses.CharacterMetadata
	//var filteredCharacters []models.Character

	url := fmt.Sprintf("https://swapi.dev/api/films/%s/", movieId)

	res, err := http.Get(url)
	if err != nil {
		return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
	}

	var movieData map[string]interface{}
	err = json.Unmarshal(body, &movieData)
	if err != nil {
		return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
	}

	charactersUrls := movieData["characters"].([]interface{})

	metadata.TotalCharacters = len(charactersUrls)

	var totalHeightCm float64 = 0

	for _, characterUrl := range charactersUrls {
		res, err := http.Get(characterUrl.(string))
		if err != nil {
			return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
		}

		var characterData map[string]interface{}
		err = json.Unmarshal(body, &characterData)
		if err != nil {
			return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
		}

		heightCm, err := strconv.ParseFloat(strings.Replace(characterData["height"].(string), ",", "", -1), 64)
		if err != nil {
			return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
		}

		totalHeightCm += heightCm

		characters = append(characters, models.Character{
			Name:   characterData["name"].(string),
			Height: characterData["height"].(string),
			Gender: characterData["gender"].(string),
		})
	}

	// Filter characters by gender, if gender filter is provided
	// if gender != "" {
	// 	totalHeightCm = 0
	// 	for _, character := range characters {
	// 		if strings.ToLower(character.Gender) == strings.ToLower(gender) {
	// 			filteredCharacters = append(filteredCharacters, character)
	// 		}
	// 	}

	// 	for _, filteredCharacters := range filteredCharacters {
	// 		heightCm, err := strconv.ParseFloat(strings.Replace(filteredCharacters.Height, ",", "", -1), 64)
	// 		if err != nil {
	// 			return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
	// 		}
	// 		totalHeightCm += heightCm
	// 	}

	// 	// Sort characters by name
	// 	sort.Slice(filteredCharacters, func(i, j int) bool {
	// 		return filteredCharacters[i].Name < filteredCharacters[j].Name
	// 	})

	// } else {
	// 	filteredCharacters = characters
	// }

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
	inches := cm * 0.3937
	feet := math.Floor(inches / 12)
	return feet
}

func cmToInch(cm float64) float64 {
	inches := cm * 0.3937
	feet := math.Floor(inches / 12)
	inches = math.Round(inches - (feet * 12))
	return inches
}