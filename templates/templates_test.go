package templates

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mehoggan/simple-wiki-web-app-go/util"
	"github.com/stretchr/testify/assert"
)

var rootPath *string = nil

func setUp() {
	_, fpath, _, _ := runtime.Caller(0)
	baseName := strings.ReplaceAll(filepath.Base(fpath), ".go", "")
	file, err := ioutil.TempDir(os.TempDir(), baseName)
	if err != nil {
		log.Fatalf("Failed to create temporary test directory in %s.",
			os.TempDir())
	}
	path, err := filepath.Abs(file)
	if err != nil {
		log.Fatalf("Failed to get absolute path of %s.", file)
	}

	if !util.Exists(path) {
		log.Fatalf("Failed to create temparary directory %s.", path)
	} else {
		log.Printf("Created temporary directory of %s.", path)
	}

	rootPath = &path
}

func tearDown() {
	log.Printf("Removing rootPath = %s...", *rootPath)
	os.RemoveAll(*rootPath)
}

func cleanString(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\t", "")
	return str
}

func generateConfigFile(rootPath string) (string, error) {
	configString := string("server:\n")
	configString += "  doc_root: \""
	configString += rootPath
	configString += "\""
	settingsFile := path.Join(rootPath, "settings.yaml")
	log.Printf("Saving %s to %s...", string(configString), settingsFile)
	err := os.WriteFile(settingsFile, []byte(configString), 0644)
	return settingsFile, err
}

func TestInstantiateTemplates(t *testing.T) {
	settingsFile, err := generateConfigFile(*rootPath)
	if err != nil {
		t.Fatalf("Failed to create settings.yml in %s.", *rootPath)
	}
	templates := InstantiateTemplates(settingsFile)
	assert.NotNilf(t, templates, "Failed to instantiate templates.")

	templatePath := path.Join(*rootPath, "view.html")
	content, err := util.LoadToString(templatePath)
	if err != nil {
		t.Fatalf("Could not load contents from %s.", templatePath)
	}
	expected := cleanString(`<h1>{{.Title}}</h1>
		<p>
			<a href="/edit/{{.Title}}">
				edit
			</a>
		</p>
		<div>
			{{printf"%s".Body}}
		</div>`)
	assert.Equalf(t, expected, cleanString(content),
		"Expected template \"\" != actual %s", content)

	templatePath = path.Join(*rootPath, "edit.html")
	content, err = util.LoadToString(templatePath)
	if err != nil {
		t.Fatalf("Could not load contents from %s.", templatePath)
	}
	expected = cleanString(`<h1>Editing {{.Title}}</h1>
			<form action="/save/{{.Title}}" method="POST">
				<div>
					<textarea name="body" rows="20" cols="80">
						{{printf "%s" .Body}}
					</textarea>
				</div>
				<div>
					<input type="submit" value="Save">
				</div>
			</form>`)
	assert.Equalf(t, expected, cleanString(content),
		"Expected template \"\" != actual %s", content)
}

func TestMain(m *testing.M) {
	log.Printf("TestMain called, running endpoint tests...")
	setUp()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}
