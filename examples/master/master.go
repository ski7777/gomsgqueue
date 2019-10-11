package master

import (
	"strconv"
	"time"

	"github.com/ski7777/gomsgqueue/pkg/interfaces"
	"github.com/ski7777/gomsgqueue/pkg/messagequeue"
)

type ExmampleMaster struct {
	mq *messagequeue.MessageQueue
}

func (m *ExmampleMaster) Run() {
	m.mq.Run()
	i := 0
	for {
		m.mq.SendMessage("Hello World "+strconv.Itoa(i), "string")
		i++
		time.Sleep(2500 * time.Millisecond)
		m.mq.SendMessageAwaitingResponse("Hello World "+strconv.Itoa(i), "string", func(pl interface{}, pt string) {
			print("Response to ")
			print(i)
			print(", Paylod:")
			if ps, ok := pl.(string); ok {
				print(ps)
			} else {
				print("---")
			}
			println()
		}, func() {})
		time.Sleep(2500 * time.Millisecond)
	}
}

func (m *ExmampleMaster) DataHandler(id string, pt string, payload interface{}, ar bool) {
	// ignoring type here
	if ar {
		if ps, ok := payload.(string); ok {
			m.mq.SendMessageResponse(ps+"...response", "string", id)
		} else {
			m.mq.SendMessageResponse(ps, pt, id)
		}
	} else {
		print("ID: ")
		print(id)
		print(", Paylod:")
		if ps, ok := payload.(string); ok {
			print(ps)
		} else {
			print("---")
		}
		println()
	}
}

func NewExampleMaster(mq *messagequeue.MessageQueue) interfaces.Master {
	m := new(ExmampleMaster)
	m.mq = mq
	m.mq.SetDataHandler(m.DataHandler)
	return m
}
