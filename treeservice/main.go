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
	roots map[int32]*actor.PID
	nextID int32
}

type HelloMsg struct {}

func (state *NodeService) Receive(context actor.Context) {
	fmt.Printf("%v\n", context.Message())
	switch msg := context.Message().(type) {
	//umbauen auf neuetree messages
	case *messages.CheckLeftMax:
		fmt.Printf("got size: %v \n", msg.MaxKey)
		context := actor.EmptyRootContext
		props := actor.PropsFromProducer(func() actor.Actor {
			return &tree.NodeActor{nil, nil, nil, nil, -1, msg.MaxKey}
		})
		pid := context.Spawn(props)

		fmt.Printf("got pid: %v \n", pid)

		state.roots[state.nextID] = pid

		fmt.Printf("new Tree with id: %v und pid: %v", state.nextID, pid )
		state.nextID++
	}
}


func main() {
	fmt.Println("Hello Tree-Service!!")

	remote.Start("localhost:8090")
	var waitgroup sync.WaitGroup

	props := actor.PropsFromProducer(
		func() actor.Actor {
			waitgroup.Add(1)
			return &NodeService{&waitgroup, make(map[int32]*actor.PID), 1001}
		})

	pid, err := actor.SpawnNamed(props, "service")
	if err == nil {
		fmt.Printf("started %v", *pid)
		waitgroup.Wait()
	} else {
		fmt.Printf("error %v", err.Error())
	}
}
