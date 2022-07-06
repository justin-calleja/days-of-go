---
title: Day 2
toc: true
summary: >
  Re-uses exported code from day 1 to implement a similar dummy server, this time using gorilla/mux and namespacing the endpoint under /api/v1.
---

## Goals

- http server serving on port `8080` with one route registered: `GET /api/v1/echo-params`. Same functionality (same code even) as in [day 1](../day-1).
- `go.mod` should `require github.com/justin-calleja/days-of-go` so `main.go` can import the package: `"github.com/justin-calleja/days-of-go/day-1/eph"`
- Attach the handler from this package to the `echo-params` route.
- Also pull in `"github.com/gorilla/mux"` and use it to set the route
- `/api/v1` should be it's own router (able to add more routes at that point).

## Walk-through

In order to import `"github.com/justin-calleja/days-of-go/day-1/eph"`, I'll start by turning `day-2/checkpoint-1` into a module. With the current working directory at `day-2/checkpoint-1`, I can:

- `go mod init github.com/justin-calleja/days-of-go/day-2/checkpoint-1`
- `go get -u github.com/justin-calleja/days-of-go`

In the `main.go` I can now import that module along with a couple of others to get to checkpoint 1 below. I can add `./content/series/days-of-go/day-2/checkpoint-1` to this repo's top level `go.work` file so I can now run checkpoint 1 from anywhere within the repo with `go run github.com/justin-calleja/days-of-go/day-2/checkpoint-1`

### Checkpoint 1

{{< code language="go" title="Checkpoint 1" id="code-checkpoint-1" expand="Show" collapse="Hide"isCollapsed="false" >}}
{{% include "/series/days-of-go/day-2/checkpoint-1/main.go" %}}{{< /code >}}

```sh
$ go run github.com/justin-calleja/days-of-go/day-2/checkpoint-1 | jq
{
  "pathParams": [
    "hello",
    "world"
  ],
  "queryParams": {
    "こんにちは": [
      "世界"
    ]
  }
}
```

This is basically using the `eph.Response` defined in [day 1](../day-1) and streaming it to `os.Stdout`.

### github.com/gorilla/mux

Add [github.com/gorilla/mux](https://github.com/gorilla/mux) with `go get -u github.com/gorilla/mux`, and create a new router and a sub router:

```go
r := mux.NewRouter()
apiV1Router := r.PathPrefix("/api/v1").Subrouter()
```

Attach the `echoParamsHandler` to the `/echo-params` route; listen on `8080`, and it's done 🙂

```go
echoParamsHandler := eph.Handler{}
apiV1Router.PathPrefix("/echo-params").Handler(&echoParamsHandler)
log.Fatal(http.ListenAndServe(":8080", r))
```

### Final version

{{< code language="go" title="main.go" nocollapse="true" collapse=" " isCollapsed="false" >}}
{{% include "/series/days-of-go/day-2/main.go" %}}{{< /code >}}

```sh
$ curl -s localhost:8080/api/v1/echo-params/hello/world\?bye\=now | jq
{
  "pathParams": [
    "api",
    "v1",
    "echo-params",
    "hello",
    "world"
  ],
  "queryParams": {
    "bye": [
      "now"
    ]
  }
}
```

Ok so the `pathParams` also includes `"api", "v1", "echo-params"` in this one. First thing that comes to mind is to extend `"github.com/justin-calleja/days-of-go/day-1/eph"`'s `Handler` struct so that it has something like a `PathPrefix` property and use that in the `TrimPrefix` call - but the goals for this exercise have already been met and there's nothing new for me to learn by fixing this 🤷
