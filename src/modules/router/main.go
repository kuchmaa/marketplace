package router

import (
	"net/http"
	"strings"
)

type TRouter struct {
	handlers map[string]func(http.ResponseWriter, *http.Request)
	static   map[string]string
	notFound func(http.ResponseWriter, *http.Request)
}

// Создание нового роутера
func Router() *TRouter {
	return &TRouter{
		handlers: make(map[string]func(http.ResponseWriter, *http.Request)),
		static:   make(map[string]string),
		notFound: func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		},
	}
}

// Регистрируем обычный обработчик
func (r *TRouter) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.handlers[pattern] = handler
}

// Регистрируем статическую директорию
func (r *TRouter) ServeStatic(route string, dir string) {
	r.static[route] = dir
}

// Запуск сервера
func (r *TRouter) RunServer(addr string) error {
	return http.ListenAndServe(addr, r)
}

func (r *TRouter) NotFound(handler func(http.ResponseWriter, *http.Request)) {
	r.notFound = func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		if handler != nil {
			handler(w, req)
		}
	}

}

// Реализация интерфейса http.Handler
func (r *TRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Статика
	for prefix, dir := range r.static {
		if strings.HasPrefix(req.URL.Path, prefix) {
			fs := http.StripPrefix(prefix, http.FileServer(http.Dir(dir)))
			fs.ServeHTTP(w, req)
			return
		}
	}

	// Хендлеры
	if handler, ok := r.handlers[req.URL.Path]; ok {
		handler(w, req)
		return
	}

	r.notFound(w, req)
}
