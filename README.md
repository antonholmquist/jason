# Jason - JSON Library for Go

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/antonholmquist/jason) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/antonholmquist/jason/master/LICENSE)

Jason intends to be an idiomatic JSON library for Go.

## Install

```
go get github.com/antonholmquist/jason`

```

## Import

```
import (
  "github.com/antonholmquist/jason"
)
```

## Sample Project

```
package main

import (
  "github.com/antonholmquist/jason"
  "log"
)

func main() {

  exampleJSON := `{
    "name": "Walter White",
    "age": 51,
    "children": [
      "junior",
      "holly"
    ],
    "other": {
      "occupation": "chemist",
      "years": 23
    }
  }`

  j, _ := jason.NewFromString(exampleJSON)

  log.Println("name:", j.Get("name").String())
  log.Println("age:", j.Get("age").Number())

  log.Println("occupation:", j.Get("other", "occupation").String())
  log.Println("years:", j.Get("other", "years").Number())

  for i, child := range j.Get("children").Array() {
    log.Printf("child %d: %s", i, child.String())
  }

}


```