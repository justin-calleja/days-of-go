---
dateCreated:
title: Day 1
toc: true
---

## Goals

- This is just for practice!
- http server serving on port `8080` with one route registered: `GET /echo-params`
- Any paths after `/echo-params` should be considered as `pathParams`.
- Response should be a JSON object with `pathParams` and `queryParams` properties, as per:
  ```go
  type response struct {
      PathParams []string `json:"pathParams"`
      QueryParams url.Values `json:"queryParams"`
  }
  ```
- For this e.g. curl request `curl -s http://localhost:8080/echo-params/hello/world\?goodbye\=world,joking\&something\=else | jq`, the response should be:
  ```json
  {
    "pathParams": ["hello", "world"],
    "queryParams": {
      "goodbye": ["world,joking"],
      "something": ["else"]
    }
  }
  ```

## Walkthrough

Put part of the provided example in the beginning section of go's [http](https://pkg.go.dev/net/http#Handler) docs:

```go
http.Handle("/echo-params", fooHandler)

log.Fatal(http.ListenAndServe(":8080", nil))
```

… in a `main` function and add the missing `fooHandler` function which… at first I thought should take as args the same args which the anon function passed to `HandleFunc` is taking in the example i.e. a `ResponseWriter` and a reference to a `Request` - but on further investigation - it takes a `struct` which implements `ServeHTTP` which is a function that takes the aforementioned types as args:

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

So, as the example in this part of the [docs](https://pkg.go.dev/net/http#Handle) show, it's something like:

```go
type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    log.Println("hello")
}

// in main:
    fooHandler := handler{}
    http.Handle("/echo-params", &fooHandler)
```

At this point, running the server with `go run main.go`, and hitting it with `curl http://localhost:8080/echo-params`, logs out `"hello"` preceded by the date / time.

Now I need to get access to path / query params. According to [this SO answer](https://stackoverflow.com/a/34315203/990159), I can use `r.URL.Path` to get access to the full path, and `strings.TrimPrefix` to remove the path the route is hosted on. I also added the same handler on both of these:

```go
http.Handle("/echo-params", &fooHandler)
http.Handle("/echo-params/", &fooHandler)
```

so both exact requests to `/echo-params` and `/echo-params/something/else` are handled by the same handler. So with this handler:

```go
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	x := strings.TrimPrefix(r.URL.Path, "/echo-params")

	if x == "" {
		x += "/"
	}

	log.Println(r.URL.Path, x)
}
```

I get this output, which is a start:

```txt
curl http://localhost:8080/echo-params
/echo-params /
curl http://localhost:8080/echo-params/
/echo-params/ /
curl http://localhost:8080/echo-params/hello
/echo-params/hello /hello
curl http://localhost:8080/echo-params/hello/world
/echo-params/hello/world /hello/world
```

At this point, I tried using: `strings.Split(path, "/")`, but that was giving me an initial empty string e.g. `["","hello","world"]` because of the first `/` in the path. So I ended up going with `func FieldsFunc(s string, f func(rune) bool) []string ` instead. The function `f` takes a `rune`, which is just and `int32`:

```go
// rune is an alias for int32 and is equivalent to int32 in all ways. It is
// used, by convention, to distinguish character values from integer values.
type rune = int32
```

Basically, using the following, I'm able to avoid the initial empty string without post-processing the "split":

```go
f := func(c rune) bool {
    return c == rune('/')
}

splitPath := strings.FieldsFunc(path, f)
```

s := strings.Split("a,b,c", ",")

    query := r.URL.Query()
    jsonByte, _ := json.Marshal(query)
    jsonString := string(jsonByte)

    log.Println(r.URL.Path, query, jsonString)
