package endpoints

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/mehoggan/simple-wiki-web-app-go/config"
	"github.com/mehoggan/simple-wiki-web-app-go/templates"
	"github.com/mehoggan/simple-wiki-web-app-go/types"
	"github.com/mehoggan/simple-wiki-web-app-go/util"
)

type Endpoints struct {
	Config     *types.Config
	Templates  *templates.Templates
	TitleRegex *regexp.Regexp
}

func (self Endpoints) getTitle(
	writter http.ResponseWriter,
	request *http.Request) (string, error) {
	match := self.TitleRegex.FindStringSubmatch(request.URL.Path)
	if match == nil {
		return "", errors.New("invalid Page Title")
	}
	return match[2], nil
}

func (self Endpoints) MakeHandler(
	fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(writter http.ResponseWriter, request *http.Request) {
		title, err := self.getTitle(writter, request)
		if err != nil {
			http.NotFound(writter, request)
			return
		}
		fn(writter, request, title)
	}
}

func (self Endpoints) ViewHandler(
	writter http.ResponseWriter,
	request *http.Request,
	title string) {
	log.Printf("Handling %s...", request.URL.Path)
	docRoot := self.Config.Server.DocRoot
	log.Printf("Going to load wiki page from %s", docRoot)
	page, err := util.Load(title, docRoot)
	if err != nil {
		writter.WriteHeader(404)
		fmt.Fprintf(writter, "<h1>Failed to find %s.txt.</h1>", title)
	} else {
		self.Templates.RenderTemplate(writter, "view", page)
	}
}

func (self Endpoints) EditHandler(
	writter http.ResponseWriter,
	request *http.Request,
	title string) {
	log.Printf("Handling %s...", request.URL.Path)
	docRoot := self.Config.Server.DocRoot
	log.Printf("Going to load wiki page from %s", docRoot)
	page, err := util.Load(title, docRoot)
	if err != nil {
		page = &types.Page{Title: title,
			Body: []byte("Please insert your text...")}
	}
	self.Templates.RenderTemplate(writter, "edit", page)
}

func (self Endpoints) SaveHandler(
	writter http.ResponseWriter,
	request *http.Request,
	title string) {
	log.Printf("Handling %s...", request.URL.Path)
	body := request.FormValue("body")
	page := &types.Page{Title: title, Body: []byte(body)}
	docRoot := self.Config.Server.DocRoot
	log.Printf("Going to save wiki page to %s", docRoot)
	err := util.Save(page, docRoot)
	if err != nil {
		http.Redirect(writter, request, "/edit/"+title,
			http.StatusInternalServerError)
	} else {
		http.Redirect(writter, request, "/edit/"+title, http.StatusFound)
	}
}

var endpoints *Endpoints
var once sync.Once

func InitializeEndpoints(configPath string) *Endpoints {
	once.Do(func() {
		config := config.Intantiate(configPath)
		templates := templates.InstantiateTemplates(configPath)
		regex := regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
		endpoints = &Endpoints{config, templates, regex}
	})
	return endpoints
}
