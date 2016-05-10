typeflow-go [![GoDoc](https://godoc.org/github.com/typeflow/typeflow-go/web?status.png)](https://godoc.org/github.com/typeflow/typeflow-go) [![Build Status](https://travis-ci.org/typeflow/typeflow-go.svg?branch=master)](https://travis-ci.org/typeflow/typeflow-go) [![Coverage Status](https://coveralls.io/repos/typeflow/typeflow-go/badge.svg?branch=master&service=github)](https://coveralls.io/github/typeflow/typeflow-go?branch=master)
=============

- [Introduction](#Introduction)
- [Quick start](#quick-start)
  - [Dependencies](#dependencies)
  - [Plain Levenshtein distance computation](#plain-levenshtein-distance-computation)
  - [Querying for similar strings](#querying-for-similar-strings)
    - [The similarity](#the-similarity)
  - [A straightforward example](#a-straightforward-example)
- [A note on similarity](#a-note-on-similarity)
- [Docs](#docs)
- [Benchmarks](#benchmarks)
- [Contribute](#contribute)
- [License](#license)

Introduction
------------
**typeflow** is a tiny package that provides a few tools around string-based searching needs.
With **typeflow** you'll be able to search for sub-string matches and get string similarity information.

Quick start
-----------
Likely you'll want to use the `WordSource` type which provides the most high-level interface in this package.

### Dependencies

This project currently only depends on:

- [prefixmap](https://github.com/alediaferia/prefixmap)

So make sure you `go get` it :)

```sh
go get github.com/alediaferia/prefixmap
```

### Plain Levenshtein distance computation
If you just need to compute the Levenshtein distance between 2 words this is what you need:

```go
typeflow.LevenshteinDistance("alessandro", "alesasndro")
```

This will return

```bash
2
```

### Querying for similar strings
This package can be used for querying against a source of strings.
For this particular need [WordSource](https://godoc.org/github.com/typeflow/typeflow-go#WordSource) has been designed specifically.

```go
ws := NewWordSource()
ws.SetSource(myListOfStrings)
matches := ws.FindMatches("query", 0.4)
```

#### The similarity
`0.4` represents the minimum similarity we are OK with. A value of 1.0 represents an exact match.

### A straightforward example

Supposedly you have a **list of words**, say country names, and you have a **partial string** which may match one or more of them according to a certain **similarity value** that is suitable for you. 
This may be the case when, for example, providing a typeahead API for populating a dropdown of suggestions (see [Google](http://google.com)).

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
"mexico",
"micronesia",
"moldova",
"monaco",
"mongolia",
"montenegro",
"morocco",
"mozambique",
"myanmar",
"namibia",
"nauru",
"nepal",
"netherlands",
"new zealand",
"nicaragua",
"niger",
"nigeria",
"norway",
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
  
  ws.SetSource(country_list)
  
  matches, err := ws.FindMatch(substr, float32(similarity))
  if err != nil {
    panic(err)
  }
  
  if len(matches) > 0 {
    fmt.Println("Found the following matches:\n")
    for _, match := range matches {
      fmt.Printf("'%s', similarity: %f\n", match.String, match.Similarity)
    }
  } else {
    fmt.Printf("No match found for '%s'.\n", substr)
  }
}

```

Output:

```bash
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
