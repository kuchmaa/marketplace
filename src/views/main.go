package views

import (
	"fmt"
	"html/template"
	"net/http"
)

type Layout interface {
	HeaderTemplate(path string)
	FooterTemplate(path string)
	AppLayout(path string)
	BodyTemplate(path string)
	Execute(w http.ResponseWriter)
}

type AppLayout struct {
	Title       string
	Description string
	appTemplate string
	CSS         map[string]string
	JS          map[string]string
	body        string
	Header      string
	Footer      string
	Data        any
}

func NewAppPage() *AppLayout {
	return &AppLayout{
		Title:       "App Page",
		Description: "This is the default app page description.",
		appTemplate: "views/layouts/app.html",
		body:        "views/pages/home.html",
		Header:      "views/components/Header.html",
		Footer:      "views/components/Footer.html",
		CSS: map[string]string{
			"main": "main.css",
		},
		JS: map[string]string{
			"index": "index.js",
		},
		Data: nil,
	}
}

func (p *AppLayout) HeaderTemplate(path string) {
	p.Header = "views/components/" + path
}

func (p *AppLayout) FooterTemplate(path string) {
	p.Footer = "views/components/" + path
}

func (p *AppLayout) AppLayout(path string) {
	p.appTemplate = "views/layouts/" + path
}

func (p *AppLayout) BodyTemplate(path string) {
	p.body = "views/pages/" + path
}

func (p *AppLayout) Execute(w http.ResponseWriter) {
	files := []string{p.appTemplate, p.body}
	if p.Header != "" {
		files = append(files, p.Header)
	}
	if p.Footer != "" {
		files = append(files, p.Footer)
	}

	// Добавляем функцию dict
	funcMap := template.FuncMap{
		"props": func(values ...interface{ any }) (map[string]interface{ any }, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict requires an even number of arguments")
			}
			dict := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				dict[key] = values[i+1]
			}
			return dict, nil
		},
	}
	tmpl := template.New("base").Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl, err = tmpl.ParseGlob("views/components/*.html")
	if err != nil {
		http.Error(w, "Component parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "app", p)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), http.StatusInternalServerError)
	}
}
