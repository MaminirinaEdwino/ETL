package cmd

import (
	"encoding/csv"
	"os"
	"strings"
	"github.com/MaminirinaEdwino/etl/src/model"
)



func NewExtractor(path string) (*model.Extractor, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	m := make(map[string]int)
	for i, name := range header {
		cleanName := strings.ToLower(strings.TrimSpace(name))
		m[cleanName] = i
	}

	return &model.Extractor{
		FilePath: path,
		Mapping:  m,
	}, nil
}

