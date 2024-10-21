package utils

import (
	"encoding/json"
	"os"
)

func SaveJsonToFile[T any](data T, filename string) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadJsonFromFile[T any](dest *T, path string) error {
	bytes, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, dest)
	if err != nil {
		return err
	}

	return nil
}
