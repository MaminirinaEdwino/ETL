package cmd

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/MaminirinaEdwino/etl/src/model"
)

func NewExtractor(path string) (*model.Extractor, []string, error) {
	var names []string
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	header, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	m := make(map[string]int)

	for i, name := range header {
		names = append(names, name)
		cleanName := strings.ToLower(strings.TrimSpace(name))
		m[cleanName] = i
	}

	return &model.Extractor{
		FilePath: path,
		Mapping:  m,
	}, names, nil
}
