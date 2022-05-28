package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Handler struct{}

type Response struct {
	PathParams  []string   `json:"pathParams"`
	QueryParams url.Values `json:"queryParams"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := Response{PathParams: []string{}, QueryParams: make(url.Values)}

	path := strings.TrimPrefix(r.URL.Path, "/echo-params")

	if path != "" && path != "/" {
		f := func(c rune) bool {
			return c == rune('/')
		}
		splitPath := strings.FieldsFunc(path, f)
		resp.PathParams = append(resp.PathParams, splitPath...)
	}

	resp.QueryParams = r.URL.Query()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	fooHandler := Handler{}
	http.Handle("/echo-params", &fooHandler)
	http.Handle("/echo-params/", &fooHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
