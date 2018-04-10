// Package shuttletracker displays the positions of RPI's shuttles on a map.
// See https://shuttles.rpi.edu and https://github.com/wtg/shuttletracker for more information.
package main

import (
	"github.com/wtg/shuttletracker/cmd/shuttletracker"
	// "github.com/wtg/shuttletracker/cmd/exporter"

)

func main() {
	shuttletracker.Run()
	// exporter.Export()
}
