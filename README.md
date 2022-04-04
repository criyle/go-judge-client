# go-judge-client

Under designing & development

[![GoDoc](https://godoc.org/github.com/criyle/go-judge-client?status.svg)](https://godoc.org/github.com/criyle/go-judge-client) [![Go Report Card](https://goreportcard.com/badge/github.com/criyle/go-judge-client)](https://goreportcard.com/report/github.com/criyle/go-judge-client) [![Release](https://img.shields.io/github/v/tag/criyle/go-judge-client)](https://github.com/criyle/go-judge-client/releases/latest)

The goal to to reimplement [syzoj/judge-v3](https://github.com/syzoj/judge-v3) in GO language using [go-sandbox](https://github.com/criyle/go-sandbox) and [go-judge](https://github.com/criyle/go-judge)

## Workflow

``` text
+-----------------------------------------------+     +--------+
| Judger (judger logic)                         | <-> | Client |
+--------------------------+-------------+------+     +--------+
| Language (Compile + Run) | ProgramConf | Data |
+--------------------------+-------------+------+
| Executor Server (RPC)    |
+--------------------------+
```

## Interfaces

- client: receive judge tasks (websocket / socket.io / RabbitMQ / REST API / gRPC stream)
- data: interface to download, cache, lock and access test data files from website (by dataId)
- language: programming language compile & execute configurations for Executor Server
- problem: parse problem definition from configuration files

## Logic

- judger: execute judge logics (compile / standard / interactive / answer submit) then collect and calculate results

## Models

- JudgeTask: judge task pushed from website (type, source, data)
- JudgeResult: judge task result send back to website
- JudgeSetting: problem setting (from yaml) and JudgeCase

## Planned API

### Client Status

Up <-> Down

### Task Progress

Client is able to report progress to the web front-end. Task should maintain its states

Planned events are:

(Preparing) -> Parsed -> (Compiling) -> Compiled ->
loop each sub task
    (Judging)
    Progressed
end
Finished

- Parsed: problem data have been downloaded and problem configuration have been parsed (pass problem config to task)
- Compiled: user code have been compiled (success / fail)
- Progressed: single test case finished (success / fail - detail message)
- Finished: all test cases finished / compile failed

Context are defined by the client to handle cancellation. Subsequent task will be cancelled and the correspond event may not be trigger.

## TODO

- [x] socket.io client with namespace
- [x] judge_v3 protocol
- [ ] executor server integration
- [ ] refactor
- [ ] syzoj problem YAML config parser
- [ ] syzoj data downloader
- [ ] syzoj compile configuration
- [ ] file io
- [ ] special judger
- [ ] interact problem
- [ ] answer submit
- [ ] demo site
- [ ] uoj support
