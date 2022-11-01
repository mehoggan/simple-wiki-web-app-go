package filesystem

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mehoggan/simple-wiki-web-app-go/types"
)

func save(page *types.Page, root string) error {
	filename := filepath.Join(root, page.Title+".txt")
	log.Printf("Writting page to (%s)...", filename)
	return os.WriteFile(filename, page.Body, 0600)
}

func load(title string, root string) (*types.Page, error) {
	filename := filepath.Join(root, title+".txt")
	log.Printf("Loading page from (%s)...", filename)
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	} else {
		return &types.Page{Title: title, Body: body}, nil
	}
}
