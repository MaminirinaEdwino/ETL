package main

import (
	"fmt"
)








func main() {
	inputFile := "road_accident_data.csv"
	outputFile := "accidents_clean.json"

	extractor, err := NewExtractor(inputFile)
	if err != nil {
		fmt.Printf("Erreur setup: %v\n", err)
		return
	}

	rawRows := make(chan []string, 100)
	transformedData := make(chan RawAccident, 100)

	go extractor.Run(rawRows)

	go func() {
		for row := range rawRows {
			acc, err := transformRow(row, extractor)
			if err != nil {
				continue 
			}

			myConfig := FilterConfig{
				MinVehicles: 5,
			}
			if shouldKeep(acc, myConfig) {
				transformedData <- acc
			}
		}
		close(transformedData) 
	}()
	err = loadToJSON(outputFile, transformedData)
	if err != nil {
		fmt.Printf("Erreur lors du chargement: %v\n", err)
	}
}
