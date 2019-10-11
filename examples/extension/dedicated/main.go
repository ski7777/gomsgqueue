package main

import (
	"github.com/ski7777/goextensioniser/pkg/extensioniser"
	"github.com/ski7777/gomsgqueue/examples/extension"
	"github.com/ski7777/gomsgqueue/pkg/messagequeue"
)

func main() {
	mq := messagequeue.NewExtension()
	mq.ConnectToMaster(extensioniser.NewDedicatedMaster())
	e := extension.NewExampleExtension(mq.GetMessageQueue())
	e.Run()
	<-make(chan int)
}
