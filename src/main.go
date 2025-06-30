package main

import (
	"database/sql"
	database "marketplace/src/database/postgres"
	"marketplace/src/database/repository"
	"marketplace/src/modules/router"
	"marketplace/src/views"
	"net/http"
	"os"
)

func main() {
	dbconfig := repository.Config{
		Host:     os.Getenv("DBHOST"),
		Port:     os.Getenv("DBPORT"),
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASS"),
		Name:     os.Getenv("DBNAME"),
	}
	db, err := database.Connect(dbconfig)
	if err != nil {
		panic(err)
	}

	r := router.Router()
	p := views.NewAppPage()

	r.HandleFunc("/", HomePage(p))
	r.HandleFunc("/about", AboutPage(p))
	r.HandleFunc("/products", ProductsPage(p))
	r.HandleFunc("/register", registerPage(p, db))
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
		type Env struct {
			User string
		}
		p.Data = Env{
			User: os.Getenv("DBUSER"),
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
func registerPage(p *views.AppLayout, db *sql.DB) router.HTTPFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:

			username := req.FormValue("username")
			email := req.FormValue("email")
			password := req.FormValue("password")

			user := &repository.User{
				Username:     username,
				Email:        email,
				PasswordHash: password,
			}

			repo := database.UserRepo(db) // создаём репозиторий

			err := repo.Create(user) // вызываем метод Create у репозитория
			if err != nil {
				http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("User created successfully!"))
			return
		default:
			p.Title = "Register"
			p.Description = "Registers to your account."
			p.BodyTemplate("register.html")
			p.Execute(w)
			return
		}
	}
}
