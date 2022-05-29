package main

import (
	"log"
	"net/http"

	"github.com/justin-calleja/days-of-go/day-1/eph"
)

func main() {
	echoParamsHandler := eph.Handler{}
	http.Handle("/echo-params", &echoParamsHandler)
	http.Handle("/echo-params/", &echoParamsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
