---
title: Day 3
toc: true
summary: >
  Use OCR to extract data from an image. Parse result in `go` to e.g. JSON or CSV.
---

## Goals

- Use OCR to extract data from an image. Parse result in `go` to e.g. JSON or CSV.

## Test driving clearlinux/tesseract-ocr

First thing's first is to test drive some OCR software. After a quick Google search, I found [tesseract-ocr](https://hub.docker.com/r/clearlinux/tesseract-ocr). In [day-3](https://github.com/justin-calleja/days-of-go/tree/main/content/series/days-of-go/day-3/checkpoint-1), I'm running:

```sh
docker run --rm --name myapp -v "$PWD":/app -w /app clearlinux/tesseract-ocr tesseract ./res/salary-info.png stdout --oem 1
```

… where `./res/salary-info.png` is just a random image I picked off of this [SO question](https://stackoverflow.com/questions/40182253/complex-headers-in-angular2-data-table).

{{< figure src="./res/salary-info.png" alt="Salary info" position="center" style="border-radius: 8px;" caption="Dummy image for OCR input: ./res/salary-info.png" captionPosition="left" captionStyle="color: black; font-weight: 700" >}}

I'm missing the data file for the English language though, so I get:

> Error opening data file /usr/share/tessdata/eng.traineddata
> Please make sure the TESSDATA_PREFIX environment variable is set to your "tessdata" directory.
> Failed loading language 'eng'
> Tesseract couldn't load any languages!
> Could not initialize tesseract.

So I need the `TESSDATA_PREFIX` env var set and I need to get the trained data for the English langauge:

> The container doesn't include any trained data for language. You need download specific language file from https://github.com/tesseract-ocr/tessdata and copy to /usr/share/tessdata, e.g to support English, you need copy eng.traineddata.
>
> **- [clearlinux/tesseract-ocr](https://hub.docker.com/r/clearlinux/tesseract-ocr)**

So I put [eng.traineddata](https://github.com/tesseract-ocr/tessdata/blob/main/eng.traineddata) in `./res`, and tried again:

```sh
docker run --rm --name myapp \
  -v "$PWD":/app \
  -e "TESSDATA_PREFIX=/app/res" \
  -w /app \
  clearlinux/tesseract-ocr tesseract \
  ./res/salary-info.png stdout --oem 1
```

This time - success - I get… pretty much the textual data depicted by the image with - as far as I can tell - no need to clean it:

```txt
HR Information Contact

Position Salary Office Extn.
Accountant $162,700 Tokyo 5407
Chief Executive Officer (CEO) $1,200,000 London 5797
Junior Technical Author $86,000 San Francisco 1562
Software Engineer $132,000 London 2558
Software Engineer $206,850 San Francisco 1314
Integration Specialist $372,000 New York 4804
Software Engineer $163,500 London 6222
Pre-Sales Support $106,450 New York 8330
Sales Assistant $145,600 New York 3990
Senior Javascript Developer $433,060 Edinburgh 6224
```

## Checkpoint 1

What I'd like to do now is:

- Try to read a file containing the `stdout` in the previous `docker` output
- If the file doesn't exist, create it as shown above (writing to file instead of stdout) and then use the same output in a `go` program to parse out the data into CSV or JSON depending on cli flags.
- If the file exists, just read from it and parse the data skipping the `docker run`.

If I needed this quickly, I'd go for a shell script to do the "file exists / run docker" part. But this is "Days of Go", so I'm going to do everything in a Go program. This gives me the opportunity to try out some more of Go's in-built modules.

So first of all, just want to get a simpler `docker` container running (simpler args with no e.g. shell variable expansions or what not):

```go
package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("docker", "run", "--rm", "grycap/cowsay")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Running it gives me something like:

```txt
 _________________________________________
/ He laughs at every joke three times...  \
| once when it's told, once when it's     |
| explained, and once when he understands |
\ it.                                     /
 -----------------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
```

… which confirms that `exec.Command` is able to pick a binary off of my `PATH` and run it. Now, I just need a way to get the present working directory since - from the docs - I know that `os/exec` won't expand env vars for me:

> Unlike the "system" library call from C and other languages, the os/exec package intentionally does not invoke the system shell and does not expand any glob patterns or handle other expansions, pipelines, or redirections typically done by shells.
>
> To expand glob patterns, either call the shell directly, taking care to escape any dangerous input, or use the path/filepath package's Glob function. To expand environment variables, use package os's ExpandEnv.
>
> **— [os/exec](https://pkg.go.dev/os/exec)**

Both of these options seem to work just fine but `os.Getwd()` feels like it should be more system independent:

```go
fmt.Println("ExpandEnv:", os.ExpandEnv("$PWD"))
path, _ := os.Getwd()
fmt.Println("Getwd:", path)
```

So with that, the first checkpoint looks something like this (which should output the same OCR output as before when running `docker` directly):

{{< code language="go" title="Checkpoint 1" id="code-checkpoint-1" expand="Show" collapse="Hide" isCollapsed="false" >}}
{{% include "/series/days-of-go/day-3/checkpoint-1/main.go" %}}{{< /code >}}

## Checkpoint 2

The easiest thing to do now is to store the whole OCR output (not that we have much in this example), in memory. I am aware that it would be more memory efficient to use streams instead but maybe I can investigate doing that in Go some other day. To do this, it's simply a matter of writing to an in-memory buffer instead of stdout:

```go
// ...
var stdout, stderr bytes.Buffer
cmd.Stdout = &stdout
cmd.Stderr = &stderr
if err := cmd.Run(); err != nil {
  log.Fatalf("OCR (via docker) failed with %s\n", err)
}

outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
if err := ioutil.WriteFile("ocr-output.txt", data, 0); err != nil {
  log.Fatalf("Writing to ocr-output.txt failed with %s\n", err)
}

fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
```


TODO:	strings.TrimSpace(...)
