package syzojclient

import (
	"net/http"
	"reflect"
	"time"

	"github.com/criyle/go-judge-client/client"
	"github.com/criyle/go-judge-client/problem"
	"github.com/criyle/go-sandbox/runner"

	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"github.com/googollee/go-socket.io/parser"
	"github.com/ugorji/go/codec"
)

const (
	buffSize  = 64
	namespace = "/judge"
)

var dialar = &engineio.Dialer{
	Transports: []transport.Transport{polling.Default, websocket.Default},
}

type ack struct {
	id uint64
}

// Client is syzoj judge client
type Client struct {
	Done chan struct{}

	token string

	socket   engineio.Conn
	tasks    chan client.Task
	progress chan []byte // msgPack encoded message
	result   chan []byte
	finish   chan client.JudgeResult
	request  chan struct{}
	ack      chan ack

	encoder *parser.Encoder
	decoder *parser.Decoder

	errCh chan error
}

// NewClient connect to socket.io endpoint
func NewClient(url, token string) (*Client, chan error, error) {
	socket, err := dialar.Dial(url, http.Header{})
	if err != nil {
		return nil, nil, err
	}

	c := &Client{
		Done:     make(chan struct{}),
		token:    token,
		socket:   socket,
		tasks:    make(chan client.Task, buffSize),
		progress: make(chan []byte, buffSize),
		result:   make(chan []byte, buffSize),
		finish:   make(chan client.JudgeResult, buffSize),
		request:  make(chan struct{}, 1),
		ack:      make(chan ack, 1),
		encoder:  parser.NewEncoder(socket),
		decoder:  parser.NewDecoder(socket),
		errCh:    make(chan error),
	}

	go c.readLoop()
	go c.writeLoop()

	return c, c.errCh, nil
}

// C c
func (c *Client) C() <-chan client.Task {
	return c.tasks
}

func (c *Client) readLoop() (err error) {
	// handle error
	defer func() {
		if err != nil {
			select {
			case c.errCh <- err:
			default:
			}
		}
	}()

	var (
		event  string
		header parser.Header
	)
	taskType := []reflect.Type{reflect.TypeOf((*parser.Buffer)(nil))}

	for {
		if err := c.decoder.DecodeHeader(&header, &event); err != nil {
			return err
		}
		switch header.Type {
		case parser.Event:
			switch event {
			case "onTask":
				// receive binary message
				args, err := c.decoder.DecodeArgs(taskType)
				if err != nil {
					return err
				}
				buf := args[0].Interface().(*parser.Buffer)

				// decode msgPack
				var task judgeTask
				if err := codec.NewDecoderBytes(buf.Data, &codec.MsgpackHandle{}).Decode(&task); err != nil {
					return err
				}
				c.tasks <- newTask(c, &task, header.ID)
			}

		case parser.Connect:
			// if connected to namespace, emit waitForTask
			if header.Namespace == namespace {
				c.request <- struct{}{}
			}

			c.decoder.DiscardLast()

		default:
			c.decoder.DiscardLast()
		}
	}
}

func (c *Client) writeLoop() (err error) {
	// handle error
	defer func() {
		if err != nil {
			select {
			case c.errCh <- err:
			default:
			}
		}
	}()

	// connect to judge
	if err := c.encoder.Encode(parser.Header{
		Type:      parser.Connect,
		Namespace: namespace,
	}, nil); err != nil {
		return err
	}

	sendProgress := func(event string, d []byte) error {
		// binary encoding
		buff := &parser.Buffer{
			Data: d,
		}

		if err := c.encoder.Encode(parser.Header{
			Type:      parser.Event,
			Namespace: namespace,
			NeedAck:   true,
		}, []interface{}{event, c.token, buff}); err != nil {
			return err
		}
		return nil
	}

	for {
		select {
		case <-c.Done:
			return

		case <-c.request:
			if err := c.encoder.Encode(parser.Header{
				Type:      parser.Event,
				Namespace: namespace,
				NeedAck:   true,
			}, []interface{}{"waitForTask", c.token}); err != nil {
				return err
			}

		case p := <-c.progress:
			if err := sendProgress("reportProgress", p); err != nil {
				return err
			}

		case r := <-c.result:
			if err := sendProgress("reportResult", r); err != nil {
				return err
			}

		case a := <-c.ack:
			if err := c.encoder.Encode(parser.Header{
				Type:      parser.Ack,
				Namespace: namespace,
				ID:        a.id,
				NeedAck:   true,
			}, nil); err != nil {
				return err
			}
		}
	}
}

type judgeTask struct {
	Content   judgeTaskContent `json:"content"`
	ExtraData string           `json:"extraData"`
}

type judgeTaskContent struct {
	TaskID   string         `json:"taskId"`
	TestData string         `json:"testData"`
	Type     int            `json:"type"`
	Priority int            `json:"priority"`
	Param    judgeParameter `json:"param"`
}

type judgeParameter struct {
	Language    string `json:"language"`
	Code        string `json:"code"`
	TimeLimit   uint64 `json:"timeLimit"`
	MemoryLimit uint64 `json:"memoryLimit"`
	// standard
	FileIOInput  *string `json:"fileIOInput"`
	FileIOOutput *string `json:"fileIOOutput"`
	// interaction
}

func newTask(c *Client, msg *judgeTask, ackID uint64) client.Task {
	task := &client.JudgeTask{
		Type: problem.Standard,
		// Code: file.SourceCode{
		// 	Language: msg.Content.Param.Language,
		// 	Code:     file.NewMemFile("src", []byte(msg.Content.Param.Code)),
		// },
		TimeLimit:   time.Duration(msg.Content.Param.TimeLimit) * time.Millisecond,
		MemoryLimit: runner.Size(msg.Content.Param.MemoryLimit << 20),
	}

	t := &Task{
		client: c,
		task:   task,
		ackID:  ackID,
		taskID: msg.Content.TaskID,

		parsed:     make(chan *problem.Config),
		compiled:   make(chan *client.ProgressCompiled),
		progressed: make(chan *client.ProgressProgressed),
		finished:   make(chan *client.JudgeResult),
	}
	go t.loop()

	return t
}
