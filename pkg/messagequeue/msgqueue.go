package messagequeue

import (
	"sync"
	"time"

	"github.com/google/uuid"
	estructs "github.com/ski7777/goextensioniser/pkg/structs"
	etypes "github.com/ski7777/goextensioniser/pkg/types"
	"github.com/ski7777/gomsgqueue/internal/structs"
	"github.com/ski7777/gomsgqueue/pkg/types"
)

const TIMEOUT = 1 * time.Minute

type MessageQueue struct {
	sendstack        chan *structs.Message
	callbacks        map[string]*structs.ResponseWait
	callbacklock     sync.Mutex
	timeout          int64
	dh               etypes.MessageHandler
	dhlock           sync.Mutex
	stopthreads      bool
	threadwait       sync.WaitGroup
	mqdh             types.MQDataHandler
	datahandlers     map[string]types.MQDataHandler
	datahandlerslock sync.Mutex
}

func (mq *MessageQueue) SendMessage(p interface{}, pt string) (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	mq.SendMessageResponse(p, pt, id.String())
	return id.String(), nil
}

func (mq *MessageQueue) SendMessageResponse(p interface{}, pt string, id string) {
	msg := new(structs.Message)
	msg.Payload = p
	msg.PayloadType = pt
	msg.ID = id
	mq.sendstack <- msg
}

func (mq *MessageQueue) SendMessageAwaitingResponse(p interface{}, pt string, rh func(interface{}, string), toh func()) (string, error) {
	msg := new(structs.Message)
	msg.Payload = p
	msg.PayloadType = pt
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	msg.ID = id.String()
	msg.ResponseHandler = rh
	msg.TimeoutHanlder = toh
	msg.AwaitingResponse = true
	mq.callbacklock.Lock()
	mq.callbacks[msg.ID] = structs.NewResponseWaitFromMessage(msg, mq.timeout)
	mq.callbacklock.Unlock()
	mq.sendstack <- msg
	return msg.ID, nil
}

func (mq *MessageQueue) SetDataHandler(dh types.MQDataHandler) {
	mq.mqdh = dh
}

func (mq *MessageQueue) Run() {
	mq.threadwait.Add(1)
	go func() {
		defer mq.threadwait.Done()
		lastcheck := time.Now().Unix()
		for !mq.stopthreads {
			if msg, ok := <-mq.sendstack; ok {
				go func() {
					mq.dhlock.Lock()
					defer mq.dhlock.Unlock()
					mq.threadwait.Add(1)
					defer mq.threadwait.Done()
					mq.dh(&estructs.Message{Payload: msg})
				}()
			}
			if lastcheck+20 < time.Now().Unix() {
				go func() {
					mq.callbacklock.Lock()
					defer mq.callbacklock.Unlock()
					mq.threadwait.Add(1)
					defer mq.threadwait.Done()
					for i, w := range mq.callbacks {
						if w.IsTimedOut() {
							go w.TimeoutHanlder()
							delete(mq.callbacks, i)
						}
					}
				}()
				lastcheck = time.Now().Unix()
			}
		}
	}()
}

func (mq *MessageQueue) Stop() {
	mq.stopthreads = true
	close(mq.sendstack)
	mq.threadwait.Wait()
}

func (mq *MessageQueue) RegisterDataHandler(t string, dh types.MQDataHandler) {
	mq.datahandlerslock.Lock()
	defer mq.datahandlerslock.Unlock()
	mq.datahandlers[t] = dh
}

func (mq *MessageQueue) UnregisterDataHandler(t string) {
	mq.datahandlerslock.Lock()
	defer mq.datahandlerslock.Unlock()
	delete(mq.datahandlers, t)
}

func (mq *MessageQueue) handleMessage(msg *structs.Message) {
	if msg.MqControl {
		if msg.PayloadType == "KeepAlive" {
			mq.handleKeepAlive(msg)
			return
		}
		if msg.PayloadType == "Abort" {
			mq.handleKeepAlive(msg)
			return
		}
	} else {
		mq.callbacklock.Lock()
		defer mq.callbacklock.Unlock()
		mq.datahandlerslock.Lock()
		defer mq.datahandlerslock.Unlock()
		if cb, ok := mq.callbacks[msg.ID]; ok {
			go cb.ResponseHandler(msg.Payload, msg.PayloadType)
			delete(mq.callbacks, msg.ID)
			return
		}
		if dh, ok := mq.datahandlers[msg.PayloadType]; ok {
			go dh(msg.ID, msg.PayloadType, msg.Payload, msg.AwaitingResponse)
			return
		}
		if mq.mqdh != nil {
			go mq.mqdh(msg.ID, msg.PayloadType, msg.Payload, msg.AwaitingResponse)
		}
	}
}

func (mq *MessageQueue) handleKeepAlive(msg *structs.Message) {
	mq.callbacklock.Lock()
	defer mq.callbacklock.Unlock()
	if v, ok := mq.callbacks[msg.ID]; ok {
		v.ResetLastPong()
	}
}

func (mq *MessageQueue) handleAbort(msg *structs.Message) {
	mq.callbacklock.Lock()
	defer mq.callbacklock.Unlock()
	delete(mq.callbacks, msg.ID)
}

func (mq *MessageQueue) handleError(e error) {}

func (mq *MessageQueue) sendMsg(msg *structs.Message) {
	m := new(estructs.Message)
	m.Payload = msg
	mq.dh(m)
}

func NewMessageQueue() *MessageQueue {
	mq := new(MessageQueue)
	mq.sendstack = make(chan *structs.Message)
	mq.callbacks = make(map[string]*structs.ResponseWait)
	mq.datahandlers = make(map[string]types.MQDataHandler)
	mq.timeout = TIMEOUT.Nanoseconds()
	return mq
}
