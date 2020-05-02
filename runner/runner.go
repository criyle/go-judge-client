package runner

import (
	"sync"

	"github.com/criyle/go-judge-client/language"
	"github.com/criyle/go-judge/pkg/envexec"
	"github.com/criyle/go-judge/pkg/pool"
)

// Runner is the task runner
type Runner struct {
	// Queue is the message queue to receive run tasks
	Queue Receiver

	// Builder builds the sandbox runner
	Builder pool.EnvBuilder

	Language language.Language

	// pool of environment to use
	pool envexec.EnvironmentPool

	// ensure init / shutdown only once
	onceInit, onceShutdown sync.Once
}

func (r *Runner) init() {
	r.pool = pool.NewPool(r.Builder)
}

// Loop status a runner in a forever loop, waiting for task and execute
// call it in new goroutine
func (r *Runner) Loop(done <-chan struct{}) {
	r.onceInit.Do(r.init)
	c := r.Queue.ReceiveC()
loop:
	for {
		select {
		case <-done:
			break loop

		case task := <-c:
			task.Done(r.run(done, task.Task()))
		}

		// check if cancel is signaled
		select {
		case <-done:
			break loop

		default:
		}
	}
	r.onceShutdown.Do(func() {
	})
}
