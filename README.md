# go-judge-client

Under designing & development

The goal to to reimplement [syzoj/judge-v3](https://github.com/syzoj/judge-v3) in GO language using [go-sandbox](https://github.com/criyle/go-sandbox) and [go-judge](https://github.com/criyle/go-judge)

## Workflow

``` text
    ^
    | Client (talk to the website)
    v
+------+    +----+
|      |<-->|Data|
|Judger|    +----+--+
|      |<-->|Problem|
+------+    +-------+
    ^
    | TaskQueue
    v
+------+   +--------+
|Runner|<->|Language|
+------+   +--------+
    ^
    | EnvExec
    v
+--------------------+
|ContainerEnvironment|
+--------------------+
```

## Interfaces

- client: receive judge tasks (websocket / socket.io / RabbitMQ / REST API)
- data: interface to download, cache, lock and access test data files from website (by dataId)
- taskqueue: message queue to send run task and receive run task result (In memory / (RabbitMQ, Redis))
- file: general file interface (disk / memory)
- language: programming language compile & execute configurations
- problem: parse problem definition from configuration files

## Judge Logic

- judger: execute judge logics (compile / standard / interactive / answer submit) and distribute as run task to queue, the collect and calculate results
- runner: receive run task and execute in sandbox environment

## Models

- JudgeTask: judge task pushed from website (type, source, data)
- JudgeResult: judge task result send back to website
- JudgeSetting: problem setting (from yaml) and JudgeCase
- RunTask: run task parameters send to run_queue
- RunResult: run task result sent back from queue

## Planned API

### Progress

Client is able to report progress to the web front-end. Task should maintain its states

Planned events are:

- Parsed: problem data have been downloaded and problem configuration have been parsed (pass problem config to task)
- Compiled: user code have been compiled (success / fail)
- Progressed: single test case finished (success / fail - detail message)
- Finished: all test cases finished / compile failed

## TODO

- [x] socket.io client with namespace
- [x] judge_v3 protocol
- [ ] syzoj problem YAML config parser
- [ ] syzoj data downloader
- [ ] syzoj compile configuration
- [ ] file io
- [ ] special judger
- [ ] interact problem
- [ ] answer submit
- [ ] demo site
