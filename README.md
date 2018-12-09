# unindexed

[![Go Report Card](https://goreportcard.com/badge/github.com/jordan-wright/unindexed)](https://goreportcard.com/report/github.com/jordan-wright/unindexed) [![GoDoc](https://godoc.org/github.com/gophish/gophish?status.svg)](https://godoc.org/github.com/gophish/gophish) ![[TravisCI](https://travis-ci.org/jordan-wright/unindexed.svg?branch=master)](https://travis-ci.org/jordan-wright/unindexed.svg?branch=master)

A Golang HTTP FileSystem that disables directory indexing.

## Motivation

By default, the `http.Dir` filesystem has [directory indexing enabled](https://www.owasp.org/index.php/OWASP_Periodic_Table_of_Vulnerabilities_-_Directory_Indexing). For example, let's say you have a `.git/` folder at the root of the folder you're serving. If someone were to request `your_url/.git/`, the contents of the folder would be listed.

This package disables directory indexing, preventing the contents from being listed.

## Installation

```
go get -u github.com/jordan-wright/unindexed
```

## Usage

The easiest way to use `unindexed` is as a drop-in replacement for `http.Dir`, which is commonly used to serve static files.

Here's a simple example using the `gorilla/mux` router:

```go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jordan-wright/unindexed"
)

func main() {
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(unindexed.Dir("../static")))
	log.Fatal(http.ListenAndServe(":8080", router))
}
```

Other examples can be found in the `examples/` directory.