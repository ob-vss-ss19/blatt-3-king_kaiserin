package main

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"sync"

	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/tree"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
)

type NodeService struct {
	waitgroup *sync.WaitGroup
	roots []*tree.NodeActor
	nextID int32
}

type HelloMsg struct {}

func (state *NodeService) Receive(context actor.Context) {
	fmt.Printf("%v\n", context.Message())
	switch msg := context.Message().(type) {
	case *messages.CheckLeftMax:
		fmt.Println("got string %v", msg.MaxKey)
	}
}


func main() {
	fmt.Println("Hello Tree-Service!!")

	remote.Start("localhost:8090")
	var waitgroup sync.WaitGroup

	props := actor.PropsFromProducer(
		func() actor.Actor {
			waitgroup.Add(1)
			return &NodeService{&waitgroup, nil, 1001}
		})

	pid, err := actor.SpawnNamed(props, "service")
	if err == nil {
		fmt.Printf("started %v", *pid)
		waitgroup.Wait()
	} else {
		fmt.Printf("error %v", err.Error())
	}
}
