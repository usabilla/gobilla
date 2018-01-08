⚠️ _**DEPRECATION WARNING**: This is an old version of the Go client for the Usabilla API. It is adviced to use [the most recent version](https://github.com/usabilla/api-go)_


# Gobilla

Gobilla is a Go client for [Usabilla API](https://usabilla.com/api).

[![Build Status](https://travis-ci.org/usabilla/gobilla.svg?branch=master)](https://travis-ci.org/usabilla/gobilla)
[![GoDoc](http://godoc.org/github.com/usabilla/gobilla?status.svg)](http://godoc.org/github.com/usabilla/gobilla)

## Getting Started

After installing Go and setting up your [GOPATH](http://golang.org/doc/code.html#GOPATH), create a `main.go` file

```go
package main

import (
    "os"
    "fmt"

    "github.com/usabilla/gobilla"
)

func main() {
    key := os.Getenv("USABILLA_API_KEY")
    secret := os.Getenv("USABILLA_API_SECRET")

    // Pass the key and secret which should be defined as ENV vars
    gb := gobilla.New(key, secret)

    b := gb.Buttons()

    // Get the first ten buttons
    params := map[string]string{"limit": "10"}
    buttons, _ := b.Get(params)

    // Print all feedback items for each button
    for _, button := range buttons.Items {
        feedback, _ := b.Feedback().Get(button.ID, nil)
        fmt.Printf("Feedback for button: %s\n%v\n", button.ID, feedback.Items)
    }
}
```

Then install Gobilla package

    go get github.com/usabilla/gobilla

Run the file

    go run main.go

And you will get all feedback items for each button.

The project includes a more detailed [example](example/main.go), which you can run from the root directory of the project

    go run example/main.go
