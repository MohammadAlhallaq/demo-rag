package ingest

import (
	"os"
	"path/filepath"
)

func LoadDocs(path string) (map[string]string, error) {
	docs := make(map[string]string)

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(p) != ".md" {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(path, p)
		docs[rel] = string(b)
		return nil
	})

	return docs, err
}
