package main

import (
	"marketplace/src/modules/router"
	"marketplace/src/views"
	"net/http"
)

func main() {
	r := router.Router()
	p := views.NewAppPage()

	r.HandleFunc("/", HomePage(p))
	r.HandleFunc("/about", AboutPage(p))
	r.HandleFunc("/products", ProductsPage(p))
	r.ServeStatic("/static/js", "./static/js")
	r.ServeStatic("/static/css", "./static/css")

	if err := r.RunServer(":3000"); err != nil {
		panic(err)
	}
}

func HomePage(p *views.AppLayout) router.HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p.Title = "Home Page"
		p.JS = map[string]string{
			"index": "index.js",
			"home":  "pages/home.js",
		}
		p.Description = "Welcome to the home page of our application."
		p.BodyTemplate("home.html")
		p.Execute(w)
	}
}

func AboutPage(p *views.AppLayout) router.HTTPFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		p.Title = "About Us"
		p.Description = "Learn more about us on this page."
		p.BodyTemplate("about.html")
		p.Execute(w)
	}

}

type ProductsData struct {
	Products map[string]ProductData
}
type ProductData struct {
	Name  string
	Price float64
}

func ProductsPage(p *views.AppLayout) router.HTTPFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		p.Title = "Products"
		p.Description = "Explore our products on this page."
		p.Data = ProductsData{
			Products: map[string]ProductData{
				"1": {Name: "Product 1", Price: 19.99},
				"2": {Name: "Product 2", Price: 29.99},
				"3": {Name: "Product 3", Price: 39.99},
			},
		}
		p.BodyTemplate("products.html")
		p.Execute(w)
	}
}
