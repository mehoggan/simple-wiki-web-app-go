package endpoints

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/mehoggan/simple-wiki-web-app-go/types"
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

func generateConfigFile() string {
	configString := string("server:\n")
	configString += "  doc_root: \""
	configString += *rootPath
	configString += "\""
	settingsFile := path.Join(*rootPath, "settings.yaml")
	log.Printf("Saving %s to %s...", string(configString), settingsFile)
	err := os.WriteFile(settingsFile, []byte(configString), 0644)
	if err != nil {
		panic(err)
	}
	return settingsFile
}

func generatePage(rootPath string, title string, t *testing.T) string {
	page := &types.Page{
		Title: title,
		Body:  []byte("This is a sample page.")}
	util.Save(page, rootPath)
	fileName := page.Title + ".txt"
	ret := path.Join(rootPath, fileName)
	if !util.Exists(ret) {
		t.Fatalf("After save %s does NOT exist.", ret)
	}
	return ret
}

func cleanString(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\t", "")
	return str
}

func TestInitializeEndpoints(t *testing.T) {
	configPath := generateConfigFile()
	_ = InitializeEndpoints(configPath)
	editTemplatePath := path.Join(*rootPath, "edit.html")
	exists := util.Exists(editTemplatePath)
	assert.Truef(t, exists,
		"Edit templates were not saved to %s after initializing endpoints.",
		editTemplatePath)
}

func TestBadEndpoint(t *testing.T) {
	var endpoints *Endpoints = InitializeEndpoints(generateConfigFile())

	req := httptest.NewRequest(http.MethodGet, "/viev/ABC", nil)
	rec := httptest.NewRecorder()
	endpoints.ViewHandler(rec, req, "view")

	res := rec.Result()
	defer res.Body.Close()
	actualByteData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Expected error to be nil got %s.", err)
	}
	actualData := cleanString(string(actualByteData))
	expectedData := cleanString("404 page not found")
	assert.Equalf(t, expectedData, actualData,
		"The response (actual) data %s != %s (expected).",
		actualData, expectedData)
	assert.Equalf(t, 404, res.StatusCode, "Expected a 404, but got a %d",
		res.StatusCode)
}

func TestMakeHandlerForViewHandlerSuccess(t *testing.T) {
	pageDataPath := generatePage(*rootPath, "ABC", t)
	defer os.Remove(pageDataPath)
	var endpoints *Endpoints = InitializeEndpoints(generateConfigFile())

	req := httptest.NewRequest(http.MethodGet, "/view/ABC", nil)
	rec := httptest.NewRecorder()
	viewHandler := endpoints.MakeHandler(endpoints.ViewHandler)
	viewHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()
	actualByteData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Expected error to be nil got %s.", err)
	}
	actualData := cleanString(string(actualByteData))
	expectedData := cleanString(`<h1>ABC</h1>
		<p>
			<ahref="/edit/ABC">
				edit
			</a>
		</p>
		<div>
			This is a sample page.
		</div>`)
	assert.Equalf(t, expectedData, actualData,
		"The response (actual) data %s != %s (expected).",
		actualData, expectedData)
	assert.Equalf(t, 200, res.StatusCode, "Expected a 200, but got a %d",
		res.StatusCode)
}

func TestViewHandlerSuccess(t *testing.T) {
	pageDataPath := generatePage(*rootPath, "ABC", t)
	defer os.Remove(pageDataPath)
	var endpoints *Endpoints = InitializeEndpoints(generateConfigFile())

	req := httptest.NewRequest(http.MethodGet, "/view/ABC", nil)
	rec := httptest.NewRecorder()
	endpoints.ViewHandler(rec, req, "view")

	res := rec.Result()
	defer res.Body.Close()
	actualByteData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Expected error to be nil got %s.", err)
	}
	actualData := cleanString(string(actualByteData))
	expectedData := cleanString(`<h1>ABC</h1>
		<p>
			<ahref="/edit/ABC">
				edit
			</a>
		</p>
		<div>
			This is a sample page.
		</div>`)
	assert.Equalf(t, expectedData, actualData,
		"The response (actual) data %s != %s (expected).",
		actualData, expectedData)
	assert.Equalf(t, 200, res.StatusCode, "Expected a 200, but got a %d",
		res.StatusCode)
}

func TestViewHandlerPageDNE(t *testing.T) {
	// We do not create the page here.
	var endpoints *Endpoints = InitializeEndpoints(generateConfigFile())

	req := httptest.NewRequest(http.MethodGet, "/view/ABC", nil)
	rec := httptest.NewRecorder()
	endpoints.ViewHandler(rec, req, "view")

	res := rec.Result()
	defer res.Body.Close()
	actualByteData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Expected error to be nil got %s.", err)
	}
	actualData := cleanString(string(actualByteData))
	expectedData := "<h1>FailedtofindABC.txt.</h1>"
	assert.Equalf(t, 404, res.StatusCode, "Expected a 404, but got a %d",
		res.StatusCode)
	assert.Equalf(t, expectedData, actualData,
		"The response (actual) data %s != %s (expected).",
		actualData, expectedData)
}

func TestEditHandlerSuccess(t *testing.T) {
	pageDataPath := generatePage(*rootPath, "ABC", t)
	defer os.Remove(pageDataPath)

	var endpoints *Endpoints = InitializeEndpoints(generateConfigFile())

	req := httptest.NewRequest(http.MethodGet, "/edit/ABC", nil)
	rec := httptest.NewRecorder()
	endpoints.EditHandler(rec, req, "edit")

	res := rec.Result()
	defer res.Body.Close()
	actualByteData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Expected error to be nil got %s.", err)
	}
	actualData := cleanString(string(actualByteData))
	expectedData := `<h1>Editing ABC</h1>
		<form action="/save/ABC" method="POST">
			<div>
				<textarea name="body" rows="20" cols="80">
					This is a sample page.
				</textarea>
			</div>
			<div>
				<input type="submit" value="Save">
			</div>
		</form>`
	expectedData = cleanString(expectedData)
	assert.Equalf(t, expectedData, actualData,
		"The response (actual) data %s != %s (expected).",
		actualData, expectedData)
	assert.Equalf(t, 200, res.StatusCode, "Expected a 200, but got a %d",
		res.StatusCode)
}

func TestEditHandlerPageDNE(t *testing.T) {
	// We do not create the page here.
	var endpoints *Endpoints = InitializeEndpoints(generateConfigFile())

	req := httptest.NewRequest(http.MethodGet, "/view/ABC", nil)
	rec := httptest.NewRecorder()
	endpoints.EditHandler(rec, req, "edit")

	res := rec.Result()
	defer res.Body.Close()
	actualByteData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Expected error to be nil got %s.", err)
	}
	actualData := cleanString(string(actualByteData))
	expectedData := cleanString(`<h1>Editing ABC</h1>
		<form action="/save/ABC" method="POST">
			<div>
				<textarea name="body" rows="20" cols="80">Please insert your text...
				</textarea>
			</div>
			<div>
				<input type="submit" value="Save">
			</div>
		</form>`)
	assert.Equalf(t, 200, res.StatusCode, "Expected a 200, but got a %d",
		res.StatusCode)
	assert.Equalf(t, expectedData, actualData,
		"The response (actual) data %s != %s (expected).",
		actualData, expectedData)
}

func TestSaveHandlerSuccess(t *testing.T) {
	pageDataPath := generatePage(*rootPath, "ABC", t)
	defer os.Remove(pageDataPath)

	var endpoints *Endpoints = InitializeEndpoints(generateConfigFile())

	req := httptest.NewRequest(http.MethodGet, "/edit/ABC", nil)
	rec := httptest.NewRecorder()
	endpoints.SaveHandler(rec, req, "edit")

	res := rec.Result()
	defer res.Body.Close()
	actualByteData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Expected error to be nil got %s.", err)
	}
	actualData := cleanString(string(actualByteData))
	expectedData := `<ahref="/edit/ABC">Found</a>.`
	expectedData = cleanString(expectedData)
	assert.Equalf(t, expectedData, actualData,
		"The response (actual) data %s != %s (expected).",
		actualData, expectedData)
	assert.Equalf(t, 302, res.StatusCode, "Expected a 200, but got a %d",
		res.StatusCode)
}

func TestMain(m *testing.M) {
	log.Printf("TestMain called, running endpoint tests...")
	setUp()
	retCode := m.Run()
	tearDown()
	os.Exit(retCode)
}
