package main

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/justin-calleja/days-of-go/day-1/eph"
)

func main() {
	resp := eph.Response{PathParams: []string{}, QueryParams: make(url.Values)}
	resp.PathParams = append(resp.PathParams, []string{"hello", "world"}...)
	resp.QueryParams.Add("こんにちは", "世界")

	json.NewEncoder(os.Stdout).Encode(resp)
}
