package main

import (
	"log"
	"os"
	"sync"

	"github.com/criyle/go-judge-client/client/syzojclient"
	"github.com/criyle/go-judge-client/judger"
	"github.com/criyle/go-judge-client/taskqueue"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	var url = os.Getenv("WEB_URL")
	c, errCh, err := syzojclient.NewClient(url, "123")
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to %s", url)

	var wg sync.WaitGroup

	done := make(chan struct{})
	q := taskqueue.NewChannelQueue(512)

	// r := &runner.Runner{
	// 	// Builder:  bu,
	// 	// Queue:    q,
	// 	// Language: &dumbLang{},
	// }
	const parallism = 4
	for i := 0; i < parallism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// r.Loop(done)
		}()
	}

	j := &judger.Judger{
		Client: c,
		Sender: q,
		// Builder: &dumbBuilder{},
	}
	go j.Loop(done)

	go func() {
		panic(<-errCh)
	}()
	<-c.Done
}

// type dumbBuilder struct {
// }

// func (b *dumbBuilder) Build([]file.File) (problem.Config, error) {
// 	const n = 100

// 	c := make([]problem.Case, 0, n)
// 	for i := 0; i < n; i++ {
// 		inputContent := strconv.Itoa(i) + " " + strconv.Itoa(i)
// 		outputContent := strconv.Itoa(i + i)
// 		c = append(c, problem.Case{
// 			// Input:  file.NewMemFile("input", []byte(inputContent)),
// 			// Answer: file.NewMemFile("output", []byte(outputContent)),
// 		})
// 	}

// 	return problem.Config{
// 		Type: "standard",
// 		Subtasks: []problem.SubTask{
// 			problem.SubTask{
// 				ScoringType: "sum",
// 				Score:       100,
// 				Cases:       c,
// 			},
// 		},
// 	}, nil
// }

// type dumbLang struct {
// }

// func (d *dumbLang) Get(name string, t language.Type) language.ExecParam {
// 	const pathEnv = "PATH=/usr/local/bin:/usr/bin:/bin"

// 	switch t {
// 	case language.TypeCompile:
// 		return language.ExecParam{
// 			Args: []string{"/usr/bin/g++", "-O2", "-o", "a", "a.cc"},
// 			Env:  []string{pathEnv},

// 			SourceFileName:    "a.cc",
// 			CompiledFileNames: []string{"a"},

// 			TimeLimit:   10 * time.Second,
// 			MemoryLimit: srunner.Size(512 << 20),
// 			ProcLimit:   100,
// 			OutputLimit: srunner.Size(64 << 10),
// 		}

// 	case language.TypeExec:
// 		return language.ExecParam{
// 			Args: []string{"a"},
// 			Env:  []string{pathEnv},

// 			SourceFileName:    "a.cc",
// 			CompiledFileNames: []string{"a"},

// 			TimeLimit:   time.Second,
// 			MemoryLimit: srunner.Size(256 << 20),
// 			ProcLimit:   1,
// 			OutputLimit: srunner.Size(64 << 10),
// 		}

// 	default:
// 		return language.ExecParam{}
// 	}
// }
