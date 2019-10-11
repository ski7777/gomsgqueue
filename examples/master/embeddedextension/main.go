package main

import (
	"github.com/ski7777/gomsgqueue/examples/extension"
	"github.com/ski7777/gomsgqueue/examples/master"
	"github.com/ski7777/gomsgqueue/pkg/messagequeue"
)

func main() {
	emq := messagequeue.NewExtension()
	mmq := messagequeue.NewMaster()
	emq.ConnectToMaster(mmq)
	e := extension.NewExampleExtension(emq.GetMessageQueue())
	m := master.NewExampleMaster(mmq.GetMessageQueue())
	go e.Run()
	m.Run()
	<-make(chan int)
}
