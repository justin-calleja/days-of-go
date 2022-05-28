package main

import (
	"log"
	"net/http"
	"strings"
)

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/echo-params")

	if path != "" && path != "/" {
		f := func(c rune) bool {
			return c == rune('/')
		}

		splitPath := strings.FieldsFunc(path, f)

		log.Println(splitPath)
	}
}

func main() {
	fooHandler := handler{}
	http.Handle("/echo-params", &fooHandler)
	http.Handle("/echo-params/", &fooHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
