package util

import (
	"errors"
	"io"
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
		log.Printf("Could not load page from %s with %s", filename, err)
		return nil, err
	} else {
		return &types.Page{Title: title, Body: body}, nil
	}
}

func LoadToString(source string) (string, error) {
	content, err := os.ReadFile(source)
	if err == nil {
		return string(content), err
	}
	return "", err
}

func Copy(source string, destination string) (int64, error) {
	sourceFile, err := os.Open(source)
	defer sourceFile.Close()
	if err != nil {
		return 0, err
	}

	// Create new file
	newFile, err := os.Create(destination)
	defer newFile.Close()
	if err != nil {
		return 0, err
	}

	return io.Copy(newFile, sourceFile)
}

func Exists(source string) bool {
	if _, err := os.Stat(source); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		log.Printf("Reporting %s does not exist with %s.", source, err)
		return false
	}
}
