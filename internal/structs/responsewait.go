package structs

import "time"

type ResponseWait struct {
	ResponseHandler func(interface{}, string)
	TimeoutHanlder  func()
	LastPong        int64
	Timeout         int64
}

func (rw *ResponseWait) IsTimedOut() bool {
	return rw.LastPong+rw.Timeout <= time.Now().UnixNano()
}

func (rw *ResponseWait) ResetLastPong() {
	rw.LastPong = time.Now().UnixNano()
}

func NewResponseWaitFromMessage(msg *Message, to int64) *ResponseWait {
	rw := new(ResponseWait)
	rw.ResponseHandler = msg.ResponseHandler
	rw.TimeoutHanlder = msg.TimeoutHanlder
	rw.LastPong = time.Now().UnixNano()
	rw.Timeout = to
	return rw
}
