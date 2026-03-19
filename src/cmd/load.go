package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	// "github.com/MaminirinaEdwino/etl/src/model"
)

func LoadToJSON(filename string, accidents <-chan map[string]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("erreur création fichier: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	// count := 0
	for acc := range accidents {		
		err := encoder.Encode(acc)

		if err != nil {
			// fmt.Printf("Erreur encodage accident %s: %v\n", acc[""], err)
			continue
		}
		
	}
	return nil
}
