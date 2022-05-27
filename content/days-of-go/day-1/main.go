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

func main() {
	fooHandler := handler{}

	http.Handle("/echo-params", &fooHandler)
	http.Handle("/echo-params/", &fooHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// m := make(url.Values)
	resp := response{PathParams: []string{}, QueryParams: make(url.Values)}

	path := strings.TrimPrefix(r.URL.Path, "/echo-params")
	if path != "" && path != "/" {
		// splitPath := strings.Split(path, "/")
		f := func(c rune) bool {
			return c == rune('/')
		}

		splitPath := strings.FieldsFunc(path, f)
		// return !unicode.IsLetter(c) && !unicode.IsNumber(c)

		resp.PathParams = append(resp.PathParams, splitPath...)
	}

	// fmt.Fprintf(w, "/")

	// log.Println(path, splitPath)
	// type Values map[string][]string

	resp.QueryParams = r.URL.Query()
	// jsonByte, _ := json.Marshal(query)
	// jsonString := string(jsonByte)

	// log.Println(r.URL.Path, query, jsonString)

	// if path == "" {
	// 	// x += "/"
	// 	// w.
	// 	fmt.Fprintf(w, "/")

	// }

	json, _ := json.Marshal(resp)
	// if err != nil {
	// 	RespondError(w, NewError(ErrorInternalError), http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)

}
