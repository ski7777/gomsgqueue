package main

import (
	"log"

	"github.com/ski7777/goextensioniser/pkg/extensioniser"
	"github.com/ski7777/gomsgqueue/examples/master"
	"github.com/ski7777/gomsgqueue/pkg/messagequeue"
)

func main() {
	cm, err := extensioniser.NewDedicatedExtension("go run $GOPATH/src/github.com/ski7777/gomsgqueue/examples/extension/dedicated/*.go")
	if err != nil {
		log.Panic(err)
	}
	mmq := messagequeue.NewMaster()
	m := master.NewExampleMaster(mmq.GetMessageQueue())
	cm.ConnectToMaster(mmq)
	m.Run()
	<-make(chan int)
}
