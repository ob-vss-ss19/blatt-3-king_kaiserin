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
	waitgroup              *sync.WaitGroup
	roots, markedForDelete map[int32]*Validation
	nextID                 int32
}

type Validation struct {
	token string
	pid   *actor.PID
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

	state.roots[state.nextID] = &Validation{tokenstring, pid}

	fmt.Printf("new Tree with id: %v und token: %v", state.nextID, tokenstring)
	context.Respond(&messages.PflanzBaumResponse{ID: state.nextID, Token: tokenstring})
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
		if val, ok := state.roots[msg.Find.ID]; ok {
			if val.token == msg.Find.Token {
				context.RequestWithCustomSender(val.pid, &messages.Insert{Key: msg.Key, Value: msg.Value}, context.Sender())
				if val, ok := state.markedForDelete[msg.Find.ID]; ok {
					if val.token == msg.Find.Token {
						fmt.Printf("removed tree mit id %v und pid %v from marked for delete", msg.Find.ID, val.pid)
						delete(state.markedForDelete, msg.Find.ID)} }
			} else {
				fmt.Println("Tree with token %v and pid %v not found!", msg.Find.Token, msg.Find.ID)
			}
		} else {
			fmt.Println("Tree with token %v and pid %v not found!", msg.Find.Token, msg.Find.ID)
		}
	case *messages.SearchCLI:
		fmt.Printf("Got Search with Key: %v, ID: %v und Token: %v \n\n",
			msg.Key, msg.Find.ID, msg.Find.Token)
		if val, ok := state.roots[msg.Find.ID]; ok {
			if val.token == msg.Find.Token {
				context.RequestWithCustomSender(val.pid, &messages.Search{Key: msg.Key}, context.Sender())
				if val, ok := state.markedForDelete[msg.Find.ID]; ok {
					if val.token == msg.Find.Token {
						fmt.Printf("removed tree mit id %v und pid %v from marked for delete", msg.Find.ID, val.pid)
						delete(state.markedForDelete, msg.Find.ID)} }
			} else {
				fmt.Println("Tree with token %v and pid %v not found!", msg.Find.Token, msg.Find.ID)
			}
		} else {
			fmt.Println("Tree with token %v and pid %v not found!", msg.Find.Token, msg.Find.ID)
		}
	case *messages.DeleteCLI:
		fmt.Printf("Got Delete with Key: %v, ID: %v und Token: %v \n\n",
			msg.Key, msg.Find.ID, msg.Find.Token)
		if val, ok := state.roots[msg.Find.ID]; ok {
			if val.token == msg.Find.Token {
				context.RequestWithCustomSender(val.pid, &messages.Delete{Key: msg.Key}, context.Sender())
				if val, ok := state.markedForDelete[msg.Find.ID]; ok {
					if val.token == msg.Find.Token {
						fmt.Printf("removed tree mit id %v und pid %v from marked for delete", msg.Find.ID, val.pid)
						delete(state.markedForDelete, msg.Find.ID)} }
			} else {
				fmt.Println("Tree with token %v and pid %v not found!", msg.Find.Token, msg.Find.ID)
			}
		} else {
			fmt.Println("Tree with token %v and pid %v not found!", msg.Find.Token, msg.Find.ID)
		}
	case *messages.TraverseCLI:
		fmt.Printf("Got Traverse with ID: %v und Token: %v \n\n", msg.Find.ID, msg.Find.Token)
		if val, ok := state.roots[msg.Find.ID]; ok {
			if val.token == msg.Find.Token {
				context.RequestWithCustomSender(val.pid, &messages.Traverse{}, context.Sender())
				if val, ok := state.markedForDelete[msg.Find.ID]; ok {
					if val.token == msg.Find.Token {
						fmt.Printf("removed tree mit id %v und pid %v from marked for delete", msg.Find.ID, val.pid)
						delete(state.markedForDelete, msg.Find.ID)} }
			} else {
				fmt.Println("Tree with token %v and pid %v not found!", msg.Find.Token, msg.Find.ID)
			}
		} else {
			fmt.Println("Tree with token %v and pid %v not found!", msg.Find.Token, msg.Find.ID)
		}
	case *messages.BaumFaellt:
		if val, ok := state.roots[msg.ID]; ok {
			if val.token == msg.Token {
				fmt.Printf("loesche tree mit id %v und pid %v", msg.ID, val.pid)
				val.pid.Stop()
				delete(state.roots, msg.ID)
			} else {
				fmt.Println("Tree with token %v and pid %v not found!", msg.Token, msg.ID)
			}
		} else {
			fmt.Println("Tree with token %v and pid %v not found!", msg.Token, msg.ID)
		}
	case *messages.DeleteTree:
		if val, ok := state.markedForDelete[msg.Delete.ID]; ok {
			if val.token == msg.Delete.Token {
				fmt.Printf("loesche tree mit id %v und pid %v", msg.Delete.ID, val.pid)
				val.pid.Stop()
				delete(state.roots, msg.Delete.ID)
				delete(state.markedForDelete, msg.Delete.ID)
			} else {
				fmt.Println("Tree with token %v and pid %v not found!", msg.Delete.Token, msg.Delete.ID)
			}
		} else {
			if val, ok := state.roots[msg.Delete.ID]; ok {
				if val.token == msg.Delete.Token {
					fmt.Printf("marked tree mit id %v und pid %v for Delete", msg.Delete.ID, val.pid)
					state.markedForDelete[msg.Delete.ID] = &Validation{msg.Delete.Token, val.pid}
				} else {
					fmt.Println("Tree with token %v and pid %v not found!", msg.Delete.Token, msg.Delete.ID)
				}
			} else {
				fmt.Println("Tree with token %v and pid %v not found!", msg.Delete.Token, msg.Delete.ID)
			}
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
			return &NodeService{&waitgroup, make(map[int32]*Validation), make(map[int32]*Validation), 1001}
		})

	pid, err := actor.SpawnNamed(props, "service")
	if err == nil {
		fmt.Printf("started %v", *pid)
		waitgroup.Wait()
	} else {
		fmt.Printf("error %v", err.Error())
	}
}
