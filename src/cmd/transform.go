package cmd

import (
	"github.com/MaminirinaEdwino/etl/src/model"
)

func TransformRow(row []string, e *model.Extractor, fieldName map[int]string) (map[string]string, error) {
	var transformData = make(map[string]string)

	for _, el := range fieldName{
		transformData[el] = e.GetValue(row, el)
	}
	return transformData, nil

}