package config

import (
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/mehoggan/simple-wiki-web-app-go/types"
	"github.com/stretchr/testify/assert"
)

func CleanString(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\t", "")
	return str
}

func TestNew(t *testing.T) {
	rootPath := t.TempDir()
	configString := string("server:\n")
	configString += string("  doc_root: \"/Users/matthew.hoggan/Desktop\"")
	settingsFile := path.Join(rootPath, "settings.yaml")
	log.Printf("Saving %v to %v...", string(configString), settingsFile)
	err := os.WriteFile(settingsFile, []byte(configString), 0644)
	if err != nil {
		panic(err)
	}
	log.Printf("Loading config from %v...", settingsFile)
	actual := Intantiate(settingsFile)
	expected := types.Config{
		Server: types.Server{DocRoot: "/Users/matthew.hoggan/Desktop"}}
	assert.Equalf(t, *actual, expected, "Configs were not equal.")
}
