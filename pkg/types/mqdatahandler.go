package types

type MQDataHandler func(string, string, interface{}, bool) // ID, Type, Payload, Awaiting Response
