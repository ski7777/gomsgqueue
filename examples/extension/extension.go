package extension

import (
	"time"

	"github.com/ski7777/gomsgqueue/pkg/interfaces"
	"github.com/ski7777/gomsgqueue/pkg/messagequeue"
)

type ExmampleExtension struct {
	mq *messagequeue.MessageQueue
}

func (e *ExmampleExtension) Run() {
	e.mq.Run()
	for {
		e.mq.SendMessage(time.Time.String(time.Now()), "string")
		time.Sleep(5 * time.Second)
	}
}

func (e *ExmampleExtension) DataHandler(id string, pt string, payload interface{}, ar bool) {
	// ignoring type here
	if ar {
		if ps, ok := payload.(string); ok {
			e.mq.SendMessageResponse(ps+"...response", "string", id)
		} else {
			e.mq.SendMessageResponse(ps, pt, id)
		}
	}
}

func NewExampleExtension(mq *messagequeue.MessageQueue) interfaces.Extension {
	e := new(ExmampleExtension)
	e.mq = mq
	e.mq.SetDataHandler(e.DataHandler)
	return e
}
