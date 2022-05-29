---
dateCreated:
title: Day 1
toc: true
summary: >
  Covers setting up an HTTP server with in-built go modules; getting access to query params and splitting path params; sending back data as JSON via `json.NewEncoder(w).Encode(resp)`.
---

## Goals

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

## Walk-through

Put part of the provided example in the beginning section of go's [http](https://pkg.go.dev/net/http#Handler) docs:

```go
http.Handle("/echo-params", fooHandler)

log.Fatal(http.ListenAndServe(":8080", nil))
```

… in a `main` function and add the missing `fooHandler` function which seems to need to be a `struct` that implements `ServeHTTP`:

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

// ...
// in main:
fooHandler := handler{}
http.Handle("/echo-params", &fooHandler)
```

### Checkpoint 1

{{< code language="go" title="Checkpoint 1" id="code-checkpoint-1" expand="Show" collapse="Hide"isCollapsed="true" >}}
{{% include "/series/days-of-go/day-1/checkpoint-1/main.go" %}}{{< /code >}}

At this point, running the server with `go run checkpoint-1/main.go`, and hitting it with `curl http://localhost:8080/echo-params`, logs out `"hello"` preceded by the date / time.

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

### Checkpoint 2

{{< code language="go" title="Checkpoint 2" id="code-checkpoint-2" expand="Show" collapse="Hide" isCollapsed="true" >}}
{{% include "/series/days-of-go/day-1/checkpoint-2/main.go" %}}{{< /code >}}

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

### Checkpoint 3

You can run this new addition by running checkpoint 3:

{{< code language="go" title="Checkpoint 3" id="code-checkpoint-3" expand="Show" collapse="Hide" isCollapsed="true" >}}
{{% include "/series/days-of-go/day-1/checkpoint-3/main.go" %}}{{< /code >}}

Moving on. The response `struct` should look like this:

```go
type response struct {
	PathParams  []string   `json:"pathParams"`
	QueryParams url.Values `json:"queryParams"`
}
```

Which should probably be created at the top of the handler with:

```go
resp := response{PathParams: []string{}, QueryParams: make(url.Values)}
```

i.e. I start off with empty arrays and `url.Values`. Then instead of logging out the path params when I have any, I can just append them to this empty array:

```go
resp.PathParams = append(resp.PathParams, splitPath...)
```

Next, I to encode `resp` as JSON, and I'll use the approach suggested [here](https://stackoverflow.com/a/37872799/990159) to avoid buffering the JSON in memory and just stream it back to the client:

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(resp)
```

### Checkpoint 4

{{< code language="go" title="Checkpoint 4" id="code-checkpoint-4" expand="Show" collapse="Hide" isCollapsed="true" >}}
{{% include "/series/days-of-go/day-1/checkpoint-4/main.go" %}}{{< /code >}}

Running checkpoint 4 should now give:

```sh
curl http://localhost:8080/echo-params
{"pathParams":[],"queryParams":{}}

curl http://localhost:8080/echo-params/
{"pathParams":[],"queryParams":{}}

curl http://localhost:8080/echo-params/hello
{"pathParams":["hello"],"queryParams":{}}

curl http://localhost:8080/echo-params/hello/world
{"pathParams":["hello","world"],"queryParams":{}}
```

All that's left are the query params, and that's as simple as:

```go
resp.QueryParams = r.URL.Query()
```

### Final version

After putting the echo params handler (that implements `http.Handler`) in it's own package in `eph/eph.go` so that it can be re-used elsewhere:

{{< code language="txt" title="go.mod" nocollapse="true" collapse=" " isCollapsed="false" >}}
{{% include "/series/days-of-go/day-1/go.mod" %}}{{< /code >}}

{{< code language="go" title="main.go" nocollapse="true" collapse=" " isCollapsed="false" >}}
{{% include "/series/days-of-go/day-1/main.go" %}}{{< /code >}}

{{< code language="go" title="eph/eph.go" nocollapse="true" collapse=" " isCollapsed="false" >}}
{{% include "/series/days-of-go/day-1/eph/eph.go" %}}{{< /code >}}

```sh
curl -s http://localhost:8080/echo-params/hello/world\?goodbye\=world,joking\&something\=else
{"pathParams":["hello","world"],"queryParams":{"goodbye":["world,joking"],"something":["else"]}}
```
