typeflow-go [![GoDoc](https://godoc.org/github.com/typeflow/typeflow-go/web?status.png)](https://godoc.org/github.com/typeflow/typeflow-go)
=============

- [Introduction](#Introduction)
- [Quick start](#quick-start)
  - [Dependencies](##dependencies)
  - [A straightforward example](##a-straightforward-example)
- [A note on similarity](#a-note-on-similarity)
- [Docs](#docs)
- [Benchmarks](#benchmarks)
- [Contribute](#contribute)
- [License](#license)

Introduction
------------
**typeflow** is a tiny package developed in my free time to learn something more about the [Levenshtein distance](https://en.wikipedia.org/wiki/Levenshtein_distance) and provide some tools around stringgo-based searching needs.

Quick start
-----------
Likely you'll want to use the `WordSource` type which provides the most high-level interface through the Levenshtein-based string matching functionalities in this package.

### Dependencies

This project currently only depends on two other projects of mine:

- [stackgo](https://github.com/alediaferia/stackgo)
- [triego](https://github.com/typeflow/triego)

So make sure you `go get` them :)

```sh
$ go get github.com/alediaferia/stackgo
$ go get github.com/typeflow/triego
```

### A straightforward example

Supposedly you have a **list of words**, say country names, and you have a **partial string** which may match one or more of them according to a certain **similarity value** that is suitable for you. This may be the case when, for example, providing a typeahead API for populating a dropdown of suggestions (see [Google](http://google.com)).

In the following example we will have a program that holds an hard-coded list of country names and accepts 2 args:

* the **substring** to look for
* the accepted minimum **similarity** for the matches

```go
package main

import (
  . "github.com/typeflow/typeflow-go"
  "strings"
  "os"
  "fmt"
  "strconv"
  "path/filepath"
)

var country_list = []string{
"Mexico",
"Micronesia",
"Moldova",
"Monaco",
"Mongolia",
"Montenegro",
"Morocco",
"Mozambique",
"Myanmar",
"Namibia",
"Nauru",
"Nepal",
"Netherlands",
"New Zealand",
"Nicaragua",
"Niger",
"Nigeria",
"Norway",
}

func printHelpAndExit() {
  fmt.Printf("usage: %s substr similarity\n", filepath.Base(os.Args[0]))
  os.Exit(0)
}

func main() {
  args := os.Args[1:]
  
  if l := len(args); l != 2 {
    fmt.Printf("Unexpected number of arguments: got %d, expected 2\n", l)
    printHelpAndExit()
  }
  
  similarity, err := strconv.ParseFloat(args[1], 32)
  if err != nil {
    fmt.Printf("Please, specify similarity as a floating point number\n")
    printHelpAndExit()
  }
  
  substr := args[0]

  // let's setup our word source
  // we will use to search for matches
  
  ws := NewWordSource()
  
  // WordSource is smart enough
  // to accept filters to apply
  // before setting up the source.
  // In this case we don't care about
  // country name case so we
  // filter the input to make sure
  // it is stored lowercase
  var filter WordFilter = func (w string) (word string, skip bool) {
    word = strings.ToLower(w)
    skip = false

    return
  }
  
  ws.SetSource(country_list, []WordFilter{ filter })
  
  matches, err := ws.FindMatch(substr, float32(similarity))
  if err != nil {
    panic(err)
  }
  
  if len(matches) > 0 {
    fmt.Println("Found the following matches:\n")
    for _, match := range matches {
      fmt.Printf("'%s', similarity: %f\n", match.Word, match.Similarity)
    }
  } else {
    fmt.Printf("No match found for '%s'.\n", substr)
  }
}

```

Output:

```
$ go run <program name> nig 0.4
```

```
Found the following matches:

'niger', similarity: 0.600000
'nigeria', similarity: 0.428571
```

A note on similarity
--------------------

The similarity between the given substring and the found match is computed using the following formula:

```math
       levenshtein(match,substr)
1.0 - ---------------------------
         max(|match|,|substr|)
```

Docs
----

I tried and will keep trying my best to keep the sources well documented.
You can help me improving the docs as well!

Docs can be found at: https://godoc.org/github.com/typeflow/typeflow-go

Benchmarks
----------

A dedicated file is [here](performance/BENCHMARKS.md) to have insights about performance
of this project. Please, help me keep it up-to-date.

Contribute
----------
I love contributions! Please create your own branch and push a merge request.

Feel free to open issues for anything :smile:

License
----------

I'm releasing this project with a MIT license included in this repository.

Copyright (c) Alessandro Diaferia <alediaferia@gmail.com>
