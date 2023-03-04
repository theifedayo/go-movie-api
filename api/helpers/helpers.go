package helpers

import (
	"math"
)

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
