package runner

import (
	"fmt"

	"github.com/criyle/go-judge-client/problem"
	"github.com/criyle/go-judge/file"
)

const maxOutput = 4 << 20 // 4M

func (r *Runner) run(done <-chan struct{}, task *RunTask) *RunTaskResult {
	switch task.Type {
	case problem.Compile:
		return r.compile(done, task.Compile)

	default:
		return r.exec(done, task.Exec)
	}
}

func getFile(files map[string]file.File, name string) ([]byte, error) {
	if f, ok := files[name]; ok {
		return f.Content()
	}
	return nil, fmt.Errorf("file %s not exists", name)
}
