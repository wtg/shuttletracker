// Package runner provides a way for programs to easily be concurrent.
//
// Anything that implements the Runnable interface can be added to a Runner,
// and then the Runner can run all of its Runnables. Runner also implements
// the Runnable interface, so they can be nested as necessary.
package runner

import (
	"sync"
)

// A Runnable has a Run method that can block forever.
type Runnable interface {
	Run()
}

// Runner tracks and runs Runnables.
type Runner struct {
	runnables []Runnable
}

// New creates a new Runner
func New() *Runner {
	runner := &Runner{runnables: []Runnable{}}
	return runner
}

// Add adds a Runnable to the Runner.
func (r *Runner) Add(runnable Runnable) {
	r.runnables = append(r.runnables, runnable)
}

// Run runs all added Runnables concurrently.
func (r *Runner) Run() {
	wg := sync.WaitGroup{}
	for i := range r.runnables {
		// We have to create a new runnable var on each iteration of the loop
		// to ensure that the goroutine we spawn runs the correct runnable.
		runnable := r.runnables[i]
		wg.Add(1)
		go func() {
			runnable.Run()
			wg.Done()
		}()
	}
	wg.Wait()
}
