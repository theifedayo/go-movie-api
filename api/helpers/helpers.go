package helpers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
