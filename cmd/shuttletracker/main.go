// Package main bundles together all of shuttletracker's subpackages
// to create, configure, and run the shuttle tracker.
package main

import (
	"github.com/wtg/shuttletracker/cmd"
)

func main() {
	cmd.Execute()
}
