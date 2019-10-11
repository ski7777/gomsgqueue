package messagequeue

import (
	"github.com/ski7777/goextensioniser/pkg/interfaces"
	estructs "github.com/ski7777/goextensioniser/pkg/structs"
	"github.com/ski7777/gomsgqueue/internal/structs"
	"github.com/ski7777/gomsgqueue/pkg/util"
)

type Extension struct {
	mq *MessageQueue
}

func (e *Extension) ConnectToMaster(m interfaces.Master) {
	m.SetDataHandler(func(msg *estructs.Message) error {
		var p *structs.Message
		var ok bool
		if p, ok = msg.Payload.(*structs.Message); !ok {
			p = new(structs.Message)
			util.GetPayload(p, msg.Payload)
		}
		go e.mq.handleMessage(p)
		return nil
	})
	m.SetErrorHandler(func(err error) {
		go e.mq.handleError(err)
	})
	e.mq.dh = m.GetDataHandler()
}

func (e *Extension) Run() {
	e.mq.Run()
}

func (e *Extension) GetMessageQueue() *MessageQueue {
	return e.mq
}

func NewExtension() *Extension {
	e := new(Extension)
	e.mq = NewMessageQueue()
	return e
}
