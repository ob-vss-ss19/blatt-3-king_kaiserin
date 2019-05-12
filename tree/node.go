package tree

import (
	"fmt"
	"sort"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
)

type SwapData struct {
	Left, Right *actor.PID
	LeftMax     int32
	Leaves      map[int32]string
}

type NodeActor struct {
	Left, Right, Parent    *actor.PID
	Leaves                 map[int32]string
	LeftMax, MaxLeaves, ID int32
	Token                  string
}

func (state *NodeActor) traverse(context actor.Context) {
	if state.Left != nil {
		leftSide, _errL := context.RequestFuture(state.Left, &messages.Traverse{}, 1*time.Second).Result()
		rightSide, _errR := context.RequestFuture(state.Right, &messages.Traverse{}, 1*time.Second).Result()
		if _errL != nil {
			println("Error with Future left happened!")
		}
		if _errR != nil {
			println("Error with Future right happened!")
		}

		lSide := leftSide.(*messages.TraverseResponse)
		rSide := rightSide.(*messages.TraverseResponse)

		lSide.Sorted = append(lSide.Sorted, rSide.Sorted...)

		context.Respond(lSide)
	} else {
		//leaves := sortKeys(state.Leaves)
		context.Respond(sortMap(state.Leaves))
	}

}

func (state *NodeActor) insert(context actor.Context) {
	msg := context.Message().(*messages.Insert)
	if state.Left != nil {
		if msg.Key > state.LeftMax {
			// msg an rechten Node, dass wert einfuegen
			context.Send(state.Right, &messages.Insert{Key: msg.Key, Value: msg.Value})
		} else {
			// msg an linken teilbaum dass er sich drum kuemmert
			context.Send(state.Left, &messages.Insert{Key: msg.Key, Value: msg.Value})
		}
	} else {
		if state.Leaves == nil {
			state.Leaves = make(map[int32]string)
		}
		// pruefen ob map schon voll
		if int32(len(state.Leaves)) == state.MaxLeaves {
			// split
			props := actor.PropsFromProducer(func() actor.Actor {
				return &NodeActor{nil, nil, context.Self(), nil, -1, state.MaxLeaves, -1, ""}
			})
			contextNew := actor.EmptyRootContext
			state.Left = contextNew.Spawn(props)
			state.Right = contextNew.Spawn(props)
			state.Leaves[msg.Key] = msg.Value
			leftMap, rightMap, leftmaximum := split(state.Leaves)
			state.LeftMax = int32(leftmaximum)
			for k, v := range leftMap {
				context.Send(state.Left, &messages.Insert{Key: k, Value: v})
			}
			for k, v := range rightMap {
				context.Send(state.Right, &messages.Insert{Key: k, Value: v})
			}
			//context.Send(state.Left, &InsertMap{leftMap})
			//context.Send(state.Right, &InsertMap{rightMap})
			state.Leaves = nil
		} else {
			state.Leaves[msg.Key] = msg.Value
		}
	}
}

func (state *NodeActor) search(context actor.Context) {
	msg := context.Message().(*messages.Search)
	if state.Left != nil {
		if msg.Key <= state.LeftMax {
			// an linken weiterschicken
			//context.RequestWithCustomSender()
			context.RequestWithCustomSender(state.Left, &messages.Search{Key: msg.Key}, context.Sender())
		} else {
			// an rechten weiterschicken
			context.RequestWithCustomSender(state.Right, &messages.Search{Key: msg.Key}, context.Sender())
		}
	} else {
		// bei mir oder gar nicht existent
		value, ok := state.Leaves[msg.Key]
		context.Respond(&messages.ScottyBeamMichHoch{Key: msg.Key, Value: value, Ok: ok})
	}
}

func (state *NodeActor) delete(context actor.Context) {
	msg := context.Message().(*messages.Delete)

	if state.Left != nil {
		if msg.Key <= state.LeftMax {
			// an linken weiterschicken
			context.RequestWithCustomSender(state.Left, &messages.Delete{Key: msg.Key}, context.Sender())
		} else {
			// an rechten weiterschicken
			context.RequestWithCustomSender(state.Right, &messages.Delete{Key: msg.Key}, context.Sender())
		}
	} else {
		// Wert in Blatt, da keinen Nachfolger mehr
		if _, ok := state.Leaves[msg.Key]; ok {
			delete(state.Leaves, msg.Key)
		}
		if len(state.Leaves) == 0 {
			// map ist leer -> actor löschen, Bruder-Actor wird zu parent
			// grandparent nicht zu parent sondern bruder von gelöschten Actor
			if state.Parent != nil {
				//  hat parent
				context.RequestWithCustomSender(state.Parent, &messages.BruderMussLos{}, context.Self())
			} else {
				//ist selbst schon root
				context.Respond(&messages.BaumFaellt{ID: state.ID, Token: state.Token})
			}

		} else {
			newMax := sortKeys(state.Leaves)
			context.Send(state.Parent, &messages.CheckLeftMax{MaxKey: int32(newMax[len(newMax)-1])})
			context.Respond(&messages.DeleteResult{Successful: true})
		}
	}
}

func (state *NodeActor) deleteChild(context actor.Context) {
	var dataToSet SwapData
	if context.Sender() == state.Left {
		// wenn sender links: rechts verknüpfen
		result, _ := context.RequestFuture(state.Right, &messages.SendMeYourData{}, 1*time.Second).Result()
		dataToSet = result.(SwapData)
	} else {
		result, _ := context.RequestFuture(state.Left, &messages.SendMeYourData{}, 1*time.Second).Result()
		dataToSet = result.(SwapData)
	}

	if dataToSet.Left != nil {
		state.Left.Stop()
		state.Right.Stop()
		state.Left = dataToSet.Left
		state.Right = dataToSet.Right

		context.Send(state.Left, &messages.SetYourPID{})
		context.Send(state.Right, &messages.SetYourPID{})

		state.LeftMax = dataToSet.LeftMax
	} else {
		state.Leaves = dataToSet.Leaves
	}
}

func (state *NodeActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Insert:
		state.insert(context)
	case *messages.Delete:
		state.delete(context)
	case *messages.Search:
		state.search(context)
	case *messages.Traverse:
		state.traverse(context)
	case *messages.ScottyBeamMichHoch:
		if msg.Ok {
			fmt.Printf("For the key '%v' there is a value '%v'! \n", msg.Key, msg.Value)
		} else {
			fmt.Printf("For the key '%v' there is NO value! \n", msg.Key)
		}
	case *messages.DeleteResult:
		if msg.Successful {
			fmt.Printf("deleting was successful! \n")
		} else {
			fmt.Printf("deleting was NOT successful! The given key does not exist.\n")
		}
	case *messages.CheckLeftMax:
		state.LeftMax = msg.MaxKey
		if state.Parent != nil {
			context.Send(state.Parent, &messages.CheckLeftMax{MaxKey: msg.MaxKey})
		}
	case *messages.BruderMussLos:
		state.deleteChild(context)
	//case *messages.IchZiehAus:
	//	state.replaceParent(context)
	//}
	case *messages.SendMeYourData:
		context.Respond(SwapData{state.Left, state.Right, state.LeftMax, state.Leaves})
		state.Right = nil
		state.Left = nil
	case *messages.SetYourPID:
		state.Parent = context.Sender()
	case *actor.Stopping:
		fmt.Printf("Stopping: %v", context.Self())
		if state.Left != nil {
			state.Right.Stop()
			state.Left.Stop()
		}
	}
}

func sortKeys(m map[int32]string) []int {
	keys := make([]int, 0)
	for k := range m {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	return keys
}

func sortMap(m map[int32]string) *messages.TraverseResponse {
	sortedKeys := sortKeys(m)
	sortedMap := make([]*messages.KeyValue, 0)

	for i := 0; i < len(m); i++ {
		sortedMap = append(sortedMap, &messages.KeyValue{Key: int32(sortedKeys[i]), Value: m[int32(sortedKeys[i])]})
	}
	return &messages.TraverseResponse{Sorted: sortedMap}
}

func split(m map[int32]string) (leftMap map[int32]string, rightMap map[int32]string, leftMax int) {
	sortedKeys := sortKeys(m)
	lengthMap := len(m) / 2
	leftMap = make(map[int32]string)
	rightMap = make(map[int32]string)

	for i := 0; i < lengthMap; i++ {
		leftMap[int32(sortedKeys[i])] = m[int32(sortedKeys[i])]
	}
	for i := lengthMap; i < len(m); i++ {
		rightMap[int32(sortedKeys[i])] = m[int32(sortedKeys[i])]
	}
	return leftMap, rightMap, sortedKeys[lengthMap-1]
}
