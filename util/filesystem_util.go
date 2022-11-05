package util

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mehoggan/simple-wiki-web-app-go/types"
)

func Save(page *types.Page, root string) error {
	filename := filepath.Join(root, page.Title+".txt")
	log.Printf("Writting page to %s...", filename)
	return os.WriteFile(filename, page.Body, 0600)
}

func Load(title string, root string) (*types.Page, error) {
	filename := filepath.Join(root, title+".txt")
	log.Printf("Loading page from %s...", filename)
	body, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("ERROR: Failed to load page from %s", filename)
		return nil, err
	} else {
		return &types.Page{Title: title, Body: body}, nil
	}
}
