package client

import (
	"time"

	"github.com/criyle/go-judge/envexec"
)

// JudgeTask contains task received from server
type JudgeTask struct {
	Type string // defines problem type
	// TestData []file.File     // test data (potential local)
	// Code     file.SourceCode // code & code language / answer submit in extra files

	// task parameters
	TimeLimit   time.Duration
	MemoryLimit envexec.Size
	Extra       map[string]interface{} // extra special parameters
}
