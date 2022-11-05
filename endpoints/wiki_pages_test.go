package endpoints

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/mehoggan/simple-wiki-web-app-go/types"
	"github.com/mehoggan/simple-wiki-web-app-go/util"
	"github.com/stretchr/testify/assert"
)

func generateConfigFile(rootPath string) string {
	configString := string("server:\n")
	configString += "  doc_root: \""
	configString += rootPath
	configString += "\""
	settingsFile := path.Join(rootPath, "settings.yaml")
	log.Printf("Saving %v to %v...", string(configString), settingsFile)
	err := os.WriteFile(settingsFile, []byte(configString), 0644)
	if err != nil {
		panic(err)
	}
	return settingsFile
}

func generatePage(rootPath string, title string) string {
	page := &types.Page{
		Title: title,
		Body:  []byte("This is a sample Page.")}
	util.Save(page, rootPath)
	fileName := page.Title + ".txt"
	return path.Join(rootPath, fileName)
}

func cleanString(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\t", "")
	return str
}

func TestViewHandlerSuccess(t *testing.T) {
	rootPath := t.TempDir()
	_ = generatePage(rootPath, "ABC")

	var endpoints *Endpoints = InitializeEndpoints(generateConfigFile(rootPath))

	req := httptest.NewRequest(http.MethodGet, "/view/ABC", nil)
	rec := httptest.NewRecorder()
	endpoints.ViewHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	actualByteData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Expected error to be nil got %v.", err)
	}
	actualData := cleanString(string(actualByteData))
	expectedData := "<h1>ABC</h1><div>ThisisasamplePage.</div>"
	assert.Equalf(t, expectedData, actualData,
		"The response (actual) data %s != %s (expected).",
		actualData, expectedData)
	assert.Equalf(t, 200, res.StatusCode, "Expected a 200, but got a %d",
		res.StatusCode)
}

func TestViewHandlerPageDNE(t *testing.T) {
	rootPath := t.TempDir()
	// We do not create the page here.

	var endpoints *Endpoints = InitializeEndpoints(generateConfigFile(rootPath))

	req := httptest.NewRequest(http.MethodGet, "/view/ABC", nil)
	rec := httptest.NewRecorder()
	endpoints.ViewHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	actualByteData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Expected error to be nil got %v.", err)
	}
	actualData := cleanString(string(actualByteData))
	expectedData := "<h1>FailedtofindABC.txt.</h1>"
	assert.Equalf(t, 404, res.StatusCode, "Expected a 404, but got a %d",
		res.StatusCode)
	assert.Equalf(t, expectedData, actualData,
		"The response (actual) data %s != %s (expected).",
		actualData, expectedData)
}
