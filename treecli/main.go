package main

import (
	"flag"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
	"sync"
)

type CLINode struct {
	waitgroup *sync.WaitGroup
}

type HelloMsg struct {}

func (state *CLINode) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *HelloMsg:
		fmt.Println("got string %v", msg)
	}
}


func main() {

	//flagCreateTree := flag.Bool("newTree", false, "creates new tree, prints out id and token")
	//flagLeafSize := flag.Int("size", 1, "size of a leaf")
	flag.Parse()

/*	var msg interface{}
	switch  {
	case *flagCreateTree:
		msg := messages.PflanzBaum(Size: int32(*flagLeafSize))

	}*/

	remote.Start("localhost:8091")
	var waitgroup sync.WaitGroup

	props := actor.PropsFromProducer(
		func() actor.Actor {
			waitgroup.Add(1)
			return &CLINode{&waitgroup}
		})

	cli := actor.Spawn(props)
	context := actor.EmptyRootContext
	remote := actor.NewPID("localhost:8090", "service")


	msg := messages.CheckLeftMax{MaxKey: 5}
	fmt.Printf("kurz vor message \n")
	context.RequestWithCustomSender(remote, &msg, cli)
	fmt.Printf("message gesendet \n")

	waitgroup.Wait()

}


