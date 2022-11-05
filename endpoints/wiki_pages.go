package endpoints

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mehoggan/simple-wiki-web-app-go/config"
	"github.com/mehoggan/simple-wiki-web-app-go/util"
)

type Endpoints struct {
	config config.Config
}

func (self Endpoints) ViewHandler(
	writter http.ResponseWriter,
	request *http.Request) {
	title := request.URL.Path[len("/view/"):]
	docRoot := self.config.Server.DocRoot
	log.Printf("Going to load wiki page from %s", docRoot)
	page, err := util.Load(title, docRoot)
	if err != nil {
		writter.WriteHeader(404)
		fmt.Fprintf(writter, "<h1>Failed to find %s.txt.</h1>", title)
	} else {
		fmt.Fprintf(writter, "<h1>%s</h1><div>%s</div>", page.Title, page.Body)
	}
}

func InitializeEndpoints(configPath string) *Endpoints {
	endpoints := &Endpoints{*config.Intantiate(configPath)}
	return endpoints
}
