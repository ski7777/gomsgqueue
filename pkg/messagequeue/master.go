package messagequeue

import (
	"github.com/ski7777/gomsgqueue/internal/structs"
	"github.com/ski7777/gomsgqueue/pkg/util"

	estructs "github.com/ski7777/goextensioniser/pkg/structs"
	etypes "github.com/ski7777/goextensioniser/pkg/types"
)

type Master struct {
	mq *MessageQueue
}

func (m *Master) SetDataHandler(dh etypes.MessageHandler) {
	m.mq.dh = dh
}

func (m *Master) SetErrorHandler(eh etypes.ErrorHandler) {}

func (m *Master) GetDataHandler() etypes.MessageHandler {
	return func(msg *estructs.Message) error {
		var p *structs.Message
		var ok bool
		if p, ok = msg.Payload.(*structs.Message); !ok {
			p = new(structs.Message)
			util.GetPayload(p, msg.Payload)
		}
		go m.mq.handleMessage(p)
		return nil
	}
}

func (m *Master) GetErrorHandler() etypes.ErrorHandler {
	return func(e error) {
		go m.mq.handleError(e)
	}
}

func (m *Master) Run() {
	m.mq.Run()
}

func (m *Master) GetMessageQueue() *MessageQueue {
	return m.mq
}

func NewMaster() *Master {
	m := new(Master)
	m.mq = NewMessageQueue()
	return m
}
