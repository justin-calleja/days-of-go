package main

import (
	"log"
	"net/http"
)

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("hello")
}

func main() {
	fooHandler := handler{}
	http.Handle("/echo-params", &fooHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
