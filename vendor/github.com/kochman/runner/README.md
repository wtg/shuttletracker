# Runner [![Build Status](https://travis-ci.org/kochman/runner.svg?branch=master)](https://travis-ci.org/kochman/runner)&nbsp;[![GoDoc](https://godoc.org/github.com/kochman/runner?status.svg)](https://godoc.org/github.com/kochman/runner)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/kochman/runner)](https://goreportcard.com/report/github.com/kochman/runner)

Runner runs anything with a `Run()` method. Perfect for API servers with background tasks, health check endpoints, or anything else that does long-lived, blocking work.

```
package main

import "github.com/kochman/runner"

func main() {
    r := runner.New()
    r.Add(myRunnable{})
    r.Add(anotherRunnable{})
    r.Run()
}
```
