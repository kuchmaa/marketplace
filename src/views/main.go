package views

import (
	"fmt"
	"html/template"
	"net/http"
)

type TAppLayout struct {
	Title       string
	Description string
	appTemplate string
	CSS         []string
	JS          []string
	body        string
	Header      string
	Footer      string
}

func NewAppPage() *TAppLayout {
	return &TAppLayout{
		Title:       "App Page",
		Description: "This is the default app page description.",
		appTemplate: "views/layouts/app.html",
		body:        "views/pages/home.html",
		Header:      "views/components/Header.html",
		Footer:      "views/components/Footer.html",
		CSS:         []string{"main.css"},
		JS:          []string{"js.js"},
	}
}

func (p *TAppLayout) HeaderTemplate(path string) {
	p.Header = "views/components/" + path
}

func (p *TAppLayout) FooterTemplate(path string) {
	p.Footer = "views/components/" + path
}

func (p *TAppLayout) AppLayout(path string) {
	p.appTemplate = "views/layout/" + path
}

func (p *TAppLayout) BodyTemplate(path string) {
	p.body = path
}

func (p *TAppLayout) Execute(w http.ResponseWriter) {
	// Собираем список шаблонов
	files := []string{p.appTemplate, p.body}
	if p.Header != "" {
		files = append(files, p.Header)
	}
	if p.Footer != "" {
		files = append(files, p.Footer)
	}

	// Добавляем функцию dict
	funcMap := template.FuncMap{
		"props": func(values ...interface{}) (map[string]interface{}, error) {
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
