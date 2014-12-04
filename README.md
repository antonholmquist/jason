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

Create value from a string. Returns an error if the string couldn't be parsed.

```go
v, err := jason.NewValueFromString(s)

```

### Create from a reader (like a http response)

Create value from a io.reader. Returns an error if the string couldn't be parsed.

```go
v, err := jason.NewValueFromReader(res.Body)

```

### Read values

Reading values is easy. If the key is invalid or type doesn't match, it will return the default value and an error.

```go
value, err := v.Get("name")
name, err := v.GetString("name")
age, err := v.GetNumber("age")
verified, err := v.GetBoolean("verified")
education, err := v.GetObject("education")
friends, err := v.GetArray("friends")

```

### Read nested values

Reading nested values is easy. If the path is invalid or type doesn't match, it will return the default value and an error.

```go
name, err := v.GetString("person", "name")
age, err := v.GetNumber("person", "age")
verified, err := v.GetBoolean("person", "verified")
education, err := v.GetObject("person", "education")
friends, err := v.GetArray("person", "friends")

```


### Loop through array

Looping through an array is done with `GetArray()`. It returns an error if the value at that keypath is null (or something else than an array).

```go

friends, err := person.GetArray("friends")
for _, friend := range friends {
  name, err := friend.GetString("name")
  age, err := friend.GetNumber("age")
}
```

### Loop through object

Looping through an object is easy. `GetObject()` returns an error if the value at that keypath is null (or something else than an object).

```go

person, err := person.GetObject("person")
for key, value := range person.Map {
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

  // Create example json
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

  // Get root value from string
  v, _ := jason.NewValueFromString(exampleJSON)

  // Read base content
  name, _ := v.GetString("name")
  age, _ := v.GetNumber("name")
  occupation, _ := v.GetString("other", "occupation")
  years, _ := v.GetNumber("other", "years")

  // Log base content
  log.Println("age:", age)
  log.Println("name:", name)
  log.Println("occupation:", occupation)
  log.Println("years:", years)

  // Loop through children array
  children, _ := v.GetArray("children")
  for i, child := range children.Slice {
    log.Printf("child %d: %s", i, child.String())
  }

  // Loop through others object
  others, _ := v.GetObject("other")
  for _, value := range others.Map() {

    s, sErr := value.AsString()
    n, nErr := value.AsNumber()

    // If it's a string, print it
    if sErr == nil {
      log.Println("string value: ", s.String())
    } 

    // If it's a number, print it
    else if nErr == nil {
      log.Println("number value: ", n.Float64())
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
