package services

func trash() {}

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"math"
// 	"net/http"
// 	"sort"
// 	"strconv"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"github.com/theifedayo/go-movie-api/api/models"
// 	"github.com/theifedayo/go-movie-api/api/responses"
// )

// func GetCharactersForMovie(movieId string, ctx *gin.Context) (int, gin.H) {
// 	var characters []models.Character
// 	var metadata responses.CharacterMetadata

// 	url := fmt.Sprintf("https://swapi.dev/api/films/%s/", movieId)

// 	res, err := http.Get(url)
// 	if err != nil {
// 		return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
// 	}

// 	defer res.Body.Close()
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
// 	}

// 	var movieData map[string]interface{}
// 	err = json.Unmarshal(body, &movieData)
// 	if err != nil {
// 		return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
// 	}

// 	charactersUrls := movieData["characters"].([]interface{})

// 	metadata.TotalCharacters = len(charactersUrls)

// 	var totalHeightCm float64 = 0

// 	for _, characterUrl := range charactersUrls {
// 		res, err := http.Get(characterUrl.(string))
// 		if err != nil {
// 			return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
// 		}

// 		defer res.Body.Close()
// 		body, err := ioutil.ReadAll(res.Body)
// 		if err != nil {
// 			return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
// 		}

// 		var characterData map[string]interface{}
// 		err = json.Unmarshal(body, &characterData)
// 		if err != nil {
// 			return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
// 		}

// 		heightCm, err := strconv.ParseFloat(strings.Replace(characterData["height"].(string), ",", "", -1), 64)
// 		if err != nil {
// 			return (http.StatusInternalServerError), gin.H{"status": "error", "data": err.Error()}
// 		}

// 		totalHeightCm += heightCm

// 		characters = append(characters, models.Character{
// 			Name:   characterData["name"].(string),
// 			Height: characterData["height"].(string),
// 			Gender: characterData["gender"].(string),
// 		})
// 	}

// 	metadata.TotalHeightCm = totalHeightCm
// 	metadata.TotalHeightFt = cmToFeet(totalHeightCm)
// 	metadata.TotalHeightIn = cmToInch(totalHeightCm)

// 	return (http.StatusOK), gin.H{
// 		"metadata": gin.H{
// 			"total_count": len(characters),
// 			"total_height": gin.H{
// 				"cm":   metadata.TotalHeightCm,
// 				"feet": metadata.TotalHeightFt,
// 				"inch": metadata.TotalHeightIn,
// 			},
// 		},
// 		"data": characters,
// 	}

// }

// func GetSortedAndFilteredCharacters(sortBy string, sortOrder string, filterByGender string) (characters []models.Character, metadata map[string]interface{}, err error) {
// 	// Get all characters from SWAPI
// 	allCharacters, err := GetAllCharacters()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// Filter characters by gender if filterByGender is provided
// 	if filterByGender != "" {
// 		filteredCharacters := make([]models.Character, 0)
// 		for _, c := range allCharacters {
// 			if c.Gender == filterByGender {
// 				filteredCharacters = append(filteredCharacters, c)
// 			}
// 		}
// 		allCharacters = filteredCharacters
// 	}

// 	// Sort characters by sortBy field and sortOrder direction
// 	switch sortBy {
// 	case "name":
// 		if sortOrder == "asc" {
// 			sort.SliceStable(allCharacters, func(i, j int) bool {
// 				return allCharacters[i].Name < allCharacters[j].Name
// 			})
// 		} else {
// 			sort.SliceStable(allCharacters, func(i, j int) bool {
// 				return allCharacters[i].Name > allCharacters[j].Name
// 			})
// 		}
// 	case "gender":
// 		if sortOrder == "asc" {
// 			sort.SliceStable(allCharacters, func(i, j int) bool {
// 				return allCharacters[i].Gender < allCharacters[j].Gender
// 			})
// 		} else {
// 			sort.SliceStable(allCharacters, func(i, j int) bool {
// 				return allCharacters[i].Gender > allCharacters[j].Gender
// 			})
// 		}
// 	case "height":
// 		if sortOrder == "asc" {
// 			sort.SliceStable(allCharacters, func(i, j int) bool {
// 				return allCharacters[i].Height < allCharacters[j].Height
// 			})
// 		} else {
// 			sort.SliceStable(allCharacters, func(i, j int) bool {
// 				return allCharacters[i].Height > allCharacters[j].Height
// 			})
// 		}
// 	}

// 	// Calculate total number of characters that match the criteria
// 	numCharacters := len(allCharacters)

// 	// Calculate total height of characters in cm and convert to feet/inches
// 	var totalHeightCm float64 = 0
// 	for _, c := range allCharacters {
// 		height, err := strconv.ParseFloat(c.Height)
// 		//heightCm, err := strconv.ParseFloat(strings.Replace(characterData["height"].(string), ",", "", -1), 64)
// 		if err != nil {
// 			continue
// 		}
// 		totalHeightCm += height
// 	}
// 	totalHeightFt := cmToFeet(totalHeightCm)
// 	totalHeightIn := cmToInch(totalHeightCm)

// 	// Create metadata map
// 	metadata = make(map[string]interface{})
// 	metadata["num_characters"] = numCharacters
// 	metadata["total_height_cm"] = totalHeightCm
// 	metadata["total_height_ft"] = totalHeightFt
// 	metadata["total_height_ft"] = totalHeightIn

// 	return allCharacters, metadata, nil

// }

// func cmToFeet(cm float64) float64 {
// 	inches := cm * 0.3937
// 	feet := math.Floor(inches / 12)
// 	return feet
// }

// func cmToInch(cm float64) float64 {
// 	inches := cm * 0.3937
// 	feet := math.Floor(inches / 12)
// 	inches = math.Round(inches - (feet * 12))
// 	return inches
// }
