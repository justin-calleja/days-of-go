package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type handler struct{}

type response struct {
	PathParams  []string   `json:"pathParams"`
	QueryParams url.Values `json:"queryParams"`
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := response{PathParams: []string{}, QueryParams: make(url.Values)}
	path := strings.TrimPrefix(r.URL.Path, "/echo-params")

	if path != "" && path != "/" {
		f := func(c rune) bool {
			return c == rune('/')
		}
		splitPath := strings.FieldsFunc(path, f)
		resp.PathParams = append(resp.PathParams, splitPath...)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	fooHandler := handler{}
	http.Handle("/echo-params", &fooHandler)
	http.Handle("/echo-params/", &fooHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
