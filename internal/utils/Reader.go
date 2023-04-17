package utils

import (
	"encoding/csv"
	"io"
	"os"
)

func ReadCSV(path string, fields int) ([][]string, error) {
	data := make([][]string, 0)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.Comma = ';'

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if len(record) == fields {
			data = append(data, record)
		}
	}
	return data, nil
}
