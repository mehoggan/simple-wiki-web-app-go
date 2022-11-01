package filesystem

import (
	"encoding/json"
	"fmt"
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
	save(page, rootPath)
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
	save(expectedPage, rootPath)
	saved := filepath.Join(rootPath, "TestPage.txt")
	actualPage, err := load(expectedTitle, rootPath)
	if err != nil {
		t.Fatalf("Failed to load (%s) while testing.", saved)
	}
	errorString := fmt.Sprintf("expected (%s) != actual (%s)",
		PageString(t, expectedPage), PageString(t, actualPage))
	assert.Equalf(t, expectedPage, actualPage, errorString)
}
