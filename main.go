package main

import (
	"fmt"

	"github.com/MaminirinaEdwino/etl/src/cmd"
	"github.com/MaminirinaEdwino/etl/src/model"
)


func main() {
	inputFile := "road_accident_data.csv"
	outputFile := "accidents_clean.json"

	extractor, err := cmd.NewExtractor(inputFile)
	if err != nil {
		fmt.Printf("Erreur setup: %v\n", err)
		return
	}

	rawRows := make(chan []string, 100)
	transformedData := make(chan model.RawAccident, 100)

	go extractor.Run(rawRows)

	go func() {
		for row := range rawRows {
			acc, err := cmd.TransformRow(row, extractor)
			if err != nil {
				continue 
			}
			myConfig := model.FilterConfig{
				MinVehicles: 12,
			}
			if cmd.ShouldKeep(acc, myConfig) {
				transformedData <- acc
			}
		}
		close(transformedData) 
	}()
	err = cmd.LoadToJSON(outputFile, transformedData)
	if err != nil {
		fmt.Printf("Erreur lors du chargement: %v\n", err)
	}
}
