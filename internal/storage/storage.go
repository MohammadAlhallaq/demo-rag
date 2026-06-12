package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"rag-demo/internal/types"
)

func SaveIndex(path string, index types.Index) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func LoadIndex(path string) (types.Index, error) {
	var index types.Index
	b, err := os.ReadFile(path)
	if err != nil {
		return index, err
	}
	err = json.Unmarshal(b, &index)
	return index, err
}
