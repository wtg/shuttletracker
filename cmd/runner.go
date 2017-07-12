package cmd

import (
	"github.com/wtg/shuttletracker/log"
	"sync"
)

type Runnable interface {
	Run()
}

type Runner struct {
	runnables []Runnable
}

func NewRunner() *Runner {
	runner := &Runner{runnables: []Runnable{}}
	return runner
}

func (r *Runner) Add(runnable Runnable) {
	r.runnables = append(r.runnables, runnable)
}

// Run runs all added runnables concurrently.
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
	log.Debug("Runnables started.")
	wg.Wait()
	log.Debug("Runnables exited.")
}
