package structs

type Message struct {
	Payload          interface{}               `json:"payload"`
	PayloadType      string                    `json:"payloadtype"`
	MqControl        bool                      `json:"mqctrl"`
	AwaitingResponse bool                      `json:"response"`
	ID               string                    `json:"ID"`
	ResponseHandler  func(interface{}, string) `json:"-"`
	TimeoutHanlder   func()                    `json:"-"`
}
