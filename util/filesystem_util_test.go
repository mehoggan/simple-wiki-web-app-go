package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/mehoggan/simple-wiki-web-app-go/types"
	"github.com/stretchr/testify/assert"
)

func PageString(t *testing.T, page *types.Page) string {
	if ret, err := json.Marshal(page); err != nil {
		t.Fatalf("Could not convert requested types.Page to string.")
	} else {
		return string(ret)
	}
	return "" // This should never be returned.
}

func TestSave(t *testing.T) {
	rootPath := t.TempDir()

	page := &types.Page{
		Title: "TestPage",
		Body:  []byte("This is a sample Page.")}
	Save(page, rootPath)
	expected := filepath.Join(rootPath, "TestPage.txt")
	_, err := os.Stat(expected)
	assert.Truef(t, err == nil, "An error occured while looking for "+expected)
}

func TestLoad(t *testing.T) {
	rootPath := t.TempDir()

	expectedTitle := "TestPage"
	expectedPage := &types.Page{
		Title: expectedTitle,
		Body:  []byte("This is a sample Page.")}
	Save(expectedPage, rootPath)
	saved := filepath.Join(rootPath, "TestPage.txt")
	actualPage, err := Load(expectedTitle, rootPath)
	if err != nil {
		t.Fatalf("Failed to load (%s) while testing.", saved)
	}
	errorString := fmt.Sprintf("expected (%s) != actual (%s)",
		PageString(t, expectedPage), PageString(t, actualPage))
	assert.Equalf(t, expectedPage, actualPage, errorString)
}

func TestLoadToString(t *testing.T) {
	rootPath := t.TempDir()

	expectedTitle := "TestPage"
	expectedPage := &types.Page{
		Title: expectedTitle,
		Body:  []byte("This is a sample Page.")}
	Save(expectedPage, rootPath)
	saved := filepath.Join(rootPath, "TestPage.txt")
	actual, err := LoadToString(saved)
	if err != nil {
		t.Fatalf("Failed in writing to %s.", saved)
	}
	expected := "This is a sample Page."
	assert.Equalf(t, expected, actual,
		"Contents of %s was not the expected of %s!!!", actual, expected)
}

func TestCopy(t *testing.T) {
	rootPath := t.TempDir()

	expectedTitle := "TestPage"
	expectedPage := &types.Page{
		Title: expectedTitle,
		Body:  []byte("This is a sample Page.")}
	Save(expectedPage, rootPath)
	saved := filepath.Join(rootPath, "TestPage.txt")
	copied := filepath.Join(rootPath, "CopyTestPage.txt")
	log.Printf("Copying %s to %s.", saved, copied)
	_, err := Copy(saved, copied)
	if err != nil {
		t.Fatalf("Failed in copying to %s.", saved)
	}
	_, err = os.Stat(copied)
	assert.False(t, errors.Is(err, os.ErrNotExist))
	savedContent, _ := LoadToString(saved)
	copiedContent, _ := LoadToString(copied)
	assert.Equalf(t, savedContent, copiedContent, "Saved %s != copied %s",
		savedContent, copiedContent)
}

func TestExist(t *testing.T) {
	rootPath := t.TempDir()
	expectedTitle := "TestPage"
	expectedPage := &types.Page{
		Title: expectedTitle,
		Body:  []byte("This is a sample Page.")}
	Save(expectedPage, rootPath)
	saved := filepath.Join(rootPath, "TestPage.txt")
	assert.Truef(t, Exists(saved), "Saved %s reported to not exist.", saved)
}

func TestDoesNotExist(t *testing.T) {
	rootPath := t.TempDir()
	// We forgot to call save
	saved := filepath.Join(rootPath, "TestPage.txt")
	assert.Falsef(t, Exists(saved), "Saved %s reported to not exist.", saved)
}
