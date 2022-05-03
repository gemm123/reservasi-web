package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gemm123/reservasi-web/internal/config"
	"github.com/gemm123/reservasi-web/internal/models"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanDate": HumanDate,
}

var app *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	app = a
}

func HumanDate(date time.Time) string {
	return date.Format("2006-01-02")
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		templates, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templates, err = templates.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templates
	}

	return myCache, nil
}

//menambahkan data di semua templates
func AddDefaultData(td *models.TemplateData, request *http.Request) *models.TemplateData {
	td.Success = app.Session.PopString(request.Context(), "success")
	td.Error = app.Session.PopString(request.Context(), "error")
	td.CSRFToken = nosurf.Token(request)
	if app.Session.Exists(request.Context(), "user_id") {
		td.IsAuthenticated = 1
	}

	return td
}

//renders template menggunakan html/template
func Template(writer http.ResponseWriter, request *http.Request, tmpl string, td *models.TemplateData) error {
	var templateCache map[string]*template.Template

	if app.UseCache {
		//mengambil template cache dari app config
		templateCache = app.TemplateCache
	} else {
		//membuat cache di setiap request
		templateCache, _ = CreateTemplateCache()
	}

	templ, ok := templateCache[tmpl]
	if !ok {
		return errors.New("cant get template from cache")
	}

	buff := new(bytes.Buffer)

	td = AddDefaultData(td, request)

	err := templ.Execute(buff, td)
	if err != nil {
		log.Fatal(err)
	}

	_, err = buff.WriteTo(writer)
	if err != nil {
		fmt.Println("error writing template to browser", err)
		return err
	}

	return nil
}
