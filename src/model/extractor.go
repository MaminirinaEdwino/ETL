package model

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Extractor struct {
	FilePath string
	Mapping  map[string]int
}

func (e *Extractor) GetValue(row []string, colName string) string {
	idx, ok := e.Mapping[strings.ToLower(colName)]
	if !ok || idx >= len(row) {
		return ""
	}
	return row[idx]
}

func (e *Extractor) Run(outChan chan<- []string) {
	file, err := os.Open(e.FilePath)
	if err != nil {
		log.Printf("Erreur ouverture fichier: %v", err)
		close(outChan)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// On ignore le header car on l'a déjà traité dans NewExtractor
	_, _ = reader.Read()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Ligne corrompue ignorée: %v", err)
			continue
		}

		// On envoie une copie de la ligne dans le channel
		line := make([]string, len(record))
		copy(line, record)
		outChan <- line
	}

	close(outChan)
	fmt.Println("Extraction terminée.")
}