package templates

import (
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"text/template"

	"github.com/mehoggan/simple-wiki-web-app-go/config"
	"github.com/mehoggan/simple-wiki-web-app-go/types"
)

type Templates struct {
	Config    *types.Config
	Templates *template.Template
}

func (self Templates) writeViewTemplateToRootDir() (int, error) {
	rootDir := self.Config.Server.DocRoot
	viewTemplatePath := path.Join(rootDir, "view.html")
	viewTemplateFile, err := os.Create(viewTemplatePath)
	if err != nil {
		log.Fatalf("Failed to create template html file %s", viewTemplatePath)
		return 0, err
	} else {
		template := `<h1>{{.Title}}</h1>
			<p>
				<a href="/edit/{{.Title}}">
					edit
				</a>
			</p>
			<div>
				{{printf "%s" .Body}}
			</div>`
		bytesCount, err := viewTemplateFile.WriteString(template)
		defer viewTemplateFile.Close()
		if err != nil {
			return 0, err
		}
		return bytesCount, err
	}
}

func (self Templates) writeEditTemplateToRootDir() (int, error) {
	rootDir := self.Config.Server.DocRoot
	editTemplatePath := path.Join(rootDir, "edit.html")
	editTemplateFile, err := os.Create(editTemplatePath)
	if err != nil {
		log.Fatalf("Failed to create template html file %s", editTemplatePath)
		return 0, err
	} else {
		template := `<h1>Editing {{.Title}}</h1>
			<form action="/save/{{.Title}}" method="POST">
				<div>
					<textarea name="body" rows="20" cols="80">
						{{printf "%s" .Body}}
					</textarea>
				</div>
				<div>
					<input type="submit" value="Save">
				</div>
			</form>`
		bytesCount, err := editTemplateFile.WriteString(template)
		defer editTemplateFile.Close()
		if err != nil {
			return 0, err
		}
		return bytesCount, err
	}
}

func (self Templates) RenderTemplate(
	writter http.ResponseWriter,
	tmpl string,
	page *types.Page) {
	err := self.Templates.ExecuteTemplate(writter, tmpl+".html", page)
	if err != nil {
		http.Error(writter, err.Error(), http.StatusInternalServerError)
	}
}

var templates *Templates = nil
var once sync.Once

func InstantiateTemplates(settingsFile string) *Templates {
	once.Do(func() {
		config := config.Intantiate(settingsFile)
		templates = &Templates{Config: config, Templates: nil}
		templates.writeEditTemplateToRootDir()
		templates.writeViewTemplateToRootDir()
		viewTemplatePath := path.Join(config.Server.DocRoot, "view.html")
		editTemplatePath := path.Join(config.Server.DocRoot, "edit.html")
		htmlTemplate := template.Must(
			template.ParseFiles(viewTemplatePath, editTemplatePath))
		templates.Templates = htmlTemplate
	})
	return templates
}
