package main

import (
	"log"
	"net/http"
	"strings"
)

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	x := strings.TrimPrefix(r.URL.Path, "/echo-params")

	if x == "" {
		x += "/"
	}

	log.Println(r.URL.Path, x)
}

func main() {
	fooHandler := handler{}
	http.Handle("/echo-params", &fooHandler)
	http.Handle("/echo-params/", &fooHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
