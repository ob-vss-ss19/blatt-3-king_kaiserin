package main

import (
	"fmt"
	"sync"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/tree"
)

type NodeService struct {
	waitgroup *sync.WaitGroup
	roots     map[int32]*actor.PID
	nextID    int32
}

func (state *NodeService) createNewTree(context actor.Context) {
	msg := context.Message().(*messages.PflanzBaum)
	fmt.Printf("got size: %v \n", msg.MaxLeaves)
	tokenstring := CreateToken(5)
	props := actor.PropsFromProducer(func() actor.Actor {
		return &tree.NodeActor{nil, nil, nil, nil, -1, msg.MaxLeaves, state.nextID, tokenstring}
	})
	pid := context.Spawn(props)

	fmt.Printf("got pid: %v \n", pid)

	state.roots[state.nextID] = pid

	fmt.Printf("new Tree with id: %v und token: %v", state.nextID, tokenstring)
	context.Respond(&messages.PflanzBaumResponse{ID:state.nextID, Token: tokenstring})
	state.nextID++
}

func (state *NodeService) Receive(context actor.Context) {
	fmt.Printf("PID Sender: %v\n\n", context.Sender())
	switch msg := context.Message().(type) {
	case *messages.PflanzBaum:
		state.createNewTree(context)
	case *messages.InsertCLI:
		fmt.Printf("Got Insert with Key: %v, Value: %v, ID: %v und Token: %v \n\n",
			msg.Key, msg.Value, msg.Find.ID, msg.Find.Token)
		pid := state.roots[msg.Find.ID]
		context.RequestWithCustomSender(pid, &messages.Insert{Key: msg.Key, Value: msg.Value}, context.Sender())
	case *messages.SearchCLI:
		fmt.Printf("Got Search with Key: %v, ID: %v und Token: %v \n\n",
			msg.Key, msg.Find.ID, msg.Find.Token)
		pid := state.roots[msg.Find.ID]
		context.RequestWithCustomSender(pid, &messages.Search{Key: msg.Key}, context.Sender())
	case *messages.DeleteCLI:
		fmt.Printf("Got Delete with Key: %v, ID: %v und Token: %v \n\n",
			msg.Key, msg.Find.ID, msg.Find.Token)
		pid := state.roots[msg.Find.ID]
		context.RequestWithCustomSender(pid, &messages.Delete{Key: msg.Key}, context.Sender())
	case *messages.TraverseCLI:
		fmt.Printf("Got Traverse with ID: %v und Token: %v \n\n", msg.Find.ID, msg.Find.Token)
		pid := state.roots[msg.Find.ID]
		context.RequestWithCustomSender(pid, &messages.Traverse{}, context.Sender())
	case *messages.BaumFaellt:
		if pid, ok := state.roots[msg.ID]; ok {
			fmt.Printf("loesche tree mit id %v und pid %v", msg.ID, pid)
			pid.Stop()
			delete(state.roots, msg.ID)
		}
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
