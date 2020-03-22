package judger

import (
	"github.com/criyle/go-judge-client/client"
	"github.com/criyle/go-judge-client/problem"
	"github.com/criyle/go-judge-client/runner"
)

// Judger receives task from client and translate to task for runner
type Judger struct {
	client.Client
	runner.Sender
	problem.Builder
}
