package runner

import (
	"sync"

	"github.com/criyle/go-judge-client/language"
)

// Runner is the task runner
type Runner struct {
	Language language.Language

	// pool of environment to use
	// pool envexec.EnvironmentPool

	// ensure init / shutdown only once
	onceInit, onceShutdown sync.Once
}

// Loop status a runner in a forever loop, waiting for task and execute
// call it in new goroutine
// func (r *Runner) Loop(done <-chan struct{}) {
// 	r.onceInit.Do(r.init)
// 	c := r.Queue.ReceiveC()
// loop:
// 	for {
// 		select {
// 		case <-done:
// 			break loop

// 		case task := <-c:
// 			task.Done(r.run(done, task.Task()))
// 		}

// 		// check if cancel is signaled
// 		select {
// 		case <-done:
// 			break loop

// 		default:
// 		}
// 	}
// 	r.onceShutdown.Do(func() {
// 	})
// }
