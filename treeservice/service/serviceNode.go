package service

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/tree"
)

type NodeService struct {
	Waitgroup              *sync.WaitGroup
	Roots, MarkedForDelete map[int32]*Validation
	NextID                 int32
}

type Validation struct {
	token string
	pid   *actor.PID
}

func (state *NodeService) createNewTree(context actor.Context) {
	msg := context.Message().(*messages.PflanzBaum)
	fmt.Printf("got size: %v \n", msg.MaxLeaves)
	tokenstring := CreateToken(5)
	nodeActor := tree.NodeActor{Left: nil, Right: nil, Parent: nil,
		Leaves: nil, LeftMax: -1, MaxLeaves: msg.MaxLeaves,
		ID: state.NextID, Token: tokenstring}
	props := actor.PropsFromProducer(func() actor.Actor {
		return &nodeActor
	})
	pid := context.Spawn(props)

	fmt.Printf("got pid: %v \n", pid)

	state.Roots[state.NextID] = &Validation{tokenstring, pid}

	if state.NextID == 1001 {
		for i := 1; i < 50; i++ {
			context.Send(pid, &messages.Insert{Key: int32(i), Value: strconv.Itoa(i)})
		}
	}
	fmt.Printf("new Tree with id: %v und token: %v", state.NextID, tokenstring)
	context.Respond(&messages.PflanzBaumResponse{ID: state.NextID, Token: tokenstring})

	state.NextID++
}

func (state *NodeService) checkIDAndToken(id int32, token string) (found bool, pid *actor.PID) {
	if val, ok := state.Roots[id]; ok {
		if val.token == token {
			return true, val.pid
		}
	}
	return false, nil
}

func (state *NodeService) removeFromMarkDelete(id int32, token string) {
	if val, ok := state.MarkedForDelete[id]; ok {
		if val.token == token {
			fmt.Printf("removed tree mit id %v und pid %v from marked for delete", id, val.pid)
			delete(state.MarkedForDelete, id)
		}
	}
}

func (state *NodeService) Receive(context actor.Context) {
	fmt.Printf("PID Sender: %v\n\n", context.Sender())
	switch msg := context.Message().(type) {
	case *messages.PflanzBaum:
		state.createNewTree(context)
	case *messages.InsertCLI:
		fmt.Printf("Got Insert with Key: %v, Value: %v, ID: %v und Token: %v \n\n",
			msg.Key, msg.Value, msg.Find.ID, msg.Find.Token)
		if ok, pid := state.checkIDAndToken(msg.Find.ID, msg.Find.Token); ok {
			context.RequestWithCustomSender(pid, &messages.Insert{Key: msg.Key, Value: msg.Value}, context.Sender())
			//context.RequestWithCustomSender(pid, &messages.Insert{Key: msg.Key + int32(1), Value: msg.Value}, context.Sender())
			context.RequestWithCustomSender(pid, &messages.Traverse{}, context.Sender())
			//context.RequestWithCustomSender(pid, &messages.Search{Key:msg.Key}, context.Sender())
			//context.RequestWithCustomSender(pid, &messages.Delete{Key:msg.Key}, context.Sender())
			state.removeFromMarkDelete(msg.Find.ID, msg.Find.Token)
		} else {
			context.Respond(&messages.TreeNotFound{NotFound: &messages.Tree{ID: msg.Find.ID, Token: msg.Find.Token}})
		}
	case *messages.SearchCLI:
		fmt.Printf("Got Search with Key: %v, ID: %v und Token: %v \n\n",
			msg.Key, msg.Find.ID, msg.Find.Token)
		if ok, pid := state.checkIDAndToken(msg.Find.ID, msg.Find.Token); ok {
			context.RequestWithCustomSender(pid, &messages.Search{Key: msg.Key}, context.Sender())
			state.removeFromMarkDelete(msg.Find.ID, msg.Find.Token)
		} else {
			context.Respond(&messages.TreeNotFound{NotFound: &messages.Tree{ID: msg.Find.ID, Token: msg.Find.Token}})
		}
	case *messages.DeleteCLI:
		fmt.Printf("Got Delete with Key: %v, ID: %v und Token: %v \n\n",
			msg.Key, msg.Find.ID, msg.Find.Token)
		if ok, pid := state.checkIDAndToken(msg.Find.ID, msg.Find.Token); ok {
			context.RequestWithCustomSender(pid, &messages.Delete{Key: msg.Key}, context.Sender())
			state.removeFromMarkDelete(msg.Find.ID, msg.Find.Token)
		} else {
			context.Respond(&messages.TreeNotFound{NotFound: &messages.Tree{ID: msg.Find.ID, Token: msg.Find.Token}})
		}
	case *messages.TraverseCLI:
		fmt.Printf("Got Traverse with ID: %v und Token: %v \n\n", msg.Find.ID, msg.Find.Token)
		if ok, pid := state.checkIDAndToken(msg.Find.ID, msg.Find.Token); ok {
			context.RequestWithCustomSender(pid, &messages.Traverse{}, context.Sender())
			state.removeFromMarkDelete(msg.Find.ID, msg.Find.Token)
		} else {
			context.Respond(&messages.TreeNotFound{NotFound: &messages.Tree{ID: msg.Find.ID, Token: msg.Find.Token}})
		}
	case *messages.BaumFaellt:
		if ok, pid := state.checkIDAndToken(msg.ID, msg.Token); ok {
			fmt.Printf("loesche tree mit id %v und pid %v", msg.ID, pid)
			pid.Stop()
			delete(state.Roots, msg.ID)
		} else {
			context.Respond(&messages.TreeNotFound{NotFound: &messages.Tree{ID: msg.ID, Token: msg.Token}})
		}
	case *messages.DeleteTree:
		if val, ok := state.MarkedForDelete[msg.Delete.ID]; ok {
			if val.token == msg.Delete.Token {
				fmt.Printf("loesche tree mit id %v und pid %v", msg.Delete.ID, val.pid)
				val.pid.Stop()
				delete(state.Roots, msg.Delete.ID)
				delete(state.MarkedForDelete, msg.Delete.ID)
				context.Respond(&messages.DeleteTreeRespond{Delete: true})
			} else {
				context.Respond(&messages.TreeNotFound{NotFound: &messages.Tree{ID: msg.Delete.ID, Token: msg.Delete.Token}})
			}
		} else {
			if val, ok := state.Roots[msg.Delete.ID]; ok {
				if val.token == msg.Delete.Token {
					fmt.Printf("marked tree mit id %v und pid %v for Delete", msg.Delete.ID, val.pid)
					state.MarkedForDelete[msg.Delete.ID] = &Validation{msg.Delete.Token, val.pid}
					context.Respond(&messages.DeleteTreeRespond{Delete: false})
				} else {
					context.Respond(&messages.TreeNotFound{NotFound: &messages.Tree{ID: msg.Delete.ID, Token: msg.Delete.Token}})
				}
			} else {
				context.Respond(&messages.TreeNotFound{NotFound: &messages.Tree{ID: msg.Delete.ID, Token: msg.Delete.Token}})
			}
		}
	}
}
