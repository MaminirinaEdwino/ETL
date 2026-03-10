package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)




func loadToJSON(filename string, accidents <-chan RawAccident) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erreur création fichier: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	count := 0
	for acc := range accidents {
		err := encoder.Encode(acc)
		if err != nil {
			fmt.Printf("Erreur encodage accident %s: %v\n", acc.Index, err)
			continue
		}

		count++
		if count%1000 == 0 {
			fmt.Printf("Chargement : %d accidents sauvegardés...\n", count)
		}
	}

	fmt.Printf("Chargement terminé : %d accidents écrits dans %s\n", count, filename)
	return nil
}

func shouldKeep(acc RawAccident, cfg FilterConfig) bool {
	if acc.Vehicles < cfg.MinVehicles {
		return false
	}
	if cfg.Severity != "" && acc.Severity != cfg.Severity {
		return false
	}
	return true
}

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
