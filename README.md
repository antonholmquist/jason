# Jason - JSON Library for Go

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/antonholmquist/jason) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/antonholmquist/jason/master/LICENSE)

Jason is an easy-to-use JSON library for Go. It's designed to be convenient for reading arbitrary JSON and to be forgiving to inconsistent content. Inspired by other libraries and improved to work well for common use cases. It currently focuses on reading JSON data rather than creating it. [API Documentation](http://godoc.org/github.com/antonholmquist/jason) can be found on godoc.org.

**Note: The API will be subject to change during 2014 if there are very good reasons to do so. On January 1st 2015 it will be frozen.**

## Data types

The following golang values are used for the JSON data types. It is consistent with how `encoding/json` uses primitive types.

- `bool`, for JSON booleans
- `float64`, for JSON numbers
- `string`, for JSON strings
- `[]*Jason`, for JSON arrays
- `map[string]*Jason`, for JSON objects
- `nil` for JSON null

## Install

```shell
go get github.com/antonholmquist/jason
```

## Import

```go
import (
  "github.com/antonholmquist/jason"
)
```

## Examples

### Create from string

Create a instance from a string. Returns an error if the string couldn't be parsed.

```go
v, err := jason.NewValueFromString(s)

```

### Create from a reader (like a http response)

Create a instance from a net/http response. Returns an error if the string couldn't be parsed.

```go
v, err := jason.NewValueFromReader(res.Body)

```

### Read values

Reading  values is easy. If the key is invalid, it will return the default value.

```go
name, err := v.Get("name").AsString()
age, err := v.Get("age").AsNumber()
verified, err := v.Get("verified").AsBoolean()
education, err := v.Get("education").AsObject()
friends,err := v.Get("friends").AsArray()

```

### Read nested values

Reading nested values is easy. If the path is invalid, it will return the default value, for instance the empty string.

```go
name, err := v.Get("person", "name").AsString()
age, err := v.Get("person", "age").AsNumber()
verified, err := v.Get("person", "verified").AsBoolean()
education, err := v.Get("person", "education").AsObject()
friends, err := v.Get("person", "friends").AsArray()

```

### Check if values exists

To check if a value exist, use `Exists()`.

```go
root.Get("person", "name").Exists()
```

### Loop through array

Looping through an array is easy and will never return an exeption. `AsArray()` returns an empty slice if the value at that keypath is null (or something else than an array).

```go

friends, err := person.Get("friends").AsArray()
for _, friend := range friends {
  name, err := friend.Get("name").AsString()
  age, err := friend.Get("age").AsNumber()
}
```

### Loop through object

Looping through an object is easy and will never return an exeption. `AsObject()` returns an empty map if the value at that keypath is null (or something else than an object).

```go

person, err := person.Get("person").AsObject()
for key, value := range person {
  ...
}
```


## Sample App

Example project demonstrating how to parse a string.

```go
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

  for key, value := range j.Get("other").Object() {
    log.Println("key: ", key)

    if value.IsString() {
      log.Println("string value: ", value.String())
    } else if value.IsNumber() {
      log.Println("number value: ", value.Number())
    }
  }

}

```

## Documentation

Documentation can be found a godoc:

https://godoc.org/github.com/antonholmquist/jason


## Test
To run the project tests:

```shell
go test
```
