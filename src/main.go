package main

import (
	"marketplace/src/modules/router"
	"marketplace/src/views"
	"net/http"
)

func main() {
	r := router.Router()

	r.HandleFunc("/", HomePage)
	r.ServeStatic("/static/js", "./static/js")
	r.ServeStatic("/static/css", "./static/css")

	if err := r.RunServer(":3000"); err != nil {
		panic(err)
	}
}
func HomePage(w http.ResponseWriter, r *http.Request) {
	p := views.NewAppPage()
	p.Execute(w)
}
