package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justin-calleja/days-of-go/day-1/eph"
)

func main() {
	r := mux.NewRouter()
	echoParamsHandler := eph.Handler{}

	apiV1Router := r.PathPrefix("/api/v1").Subrouter()
	apiV1Router.PathPrefix("/echo-params").Handler(&echoParamsHandler)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
