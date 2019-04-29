package tree

import (
	"fmt"
	"sort"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
)

type NodeActor struct {
	Left      *actor.PID
	LeftMax   int32
	Right     *actor.PID
	Leaves    map[int32]string
	MaxLeaves int32
	Parent    *actor.PID
}

func (state *NodeActor) traverse(context actor.Context) {
	if state.Left != nil {
		leftSide, _err := context.RequestFuture(state.Left, &messages.Traverse{}, 1*time.Second).Result()
		rightSide, _err := context.RequestFuture(state.Right, &messages.Traverse{}, 1*time.Second).Result()
		if _err != nil {
			println("Error with Futures happened!")
		}

		lSide := leftSide.(map[int32]string)
		rSide := rightSide.(map[int32]string)

		for k, v := range rSide {
			lSide[k] = v
		}
		/*		var full map[int32]string
				full = append(full, *lSide...)
				full = append(full, *rSide...)*/

		if state.Parent != nil {
			context.Respond(lSide)
		} else {
			fmt.Printf("All keys in Tree sorted: %v\n", lSide)
		}
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
				return &NodeActor{nil, -1, nil, nil, 4, context.Self()}
			})
			state.Left = context.Spawn(props)
			state.Right = context.Spawn(props)
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
		if value, ok := state.Leaves[msg.Key]; ok {
			context.Respond(&messages.ScottyBeamMichHoch{Key: msg.Key, Value: value, Ok: ok})
		} else {
			context.Respond(&messages.ScottyBeamMichHoch{Key: msg.Key, Value: value, Ok: ok})
		}
	}
}

func (state *NodeActor) delete(context actor.Context) {
	msg := context.Message().(*messages.Delete)
	//TODO search if Blatt vorhanden

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
			context.Send(state.Parent, &messages.BruderMussLos{})
		} else {
			newMax := sortKeys(state.Leaves)
			context.Send(state.Parent, &messages.CheckLeftMax{MaxKey: int32(newMax[len(newMax)-1])})
			context.Respond(&messages.DeleteResult{Successful: true})
		}
	}
}

func (state *NodeActor) deleteChild(context actor.Context) {
	if context.Sender() == state.Left {
		// wenn sender links: rechts verknüpfen
		context.Sender().Stop()
		context.RequestWithCustomSender(state.Parent, &messages.IchZiehAus{MyMax: state.LeftMax}, state.Right)
	} else {
		context.Sender().Stop()
		context.RequestWithCustomSender(state.Parent, &messages.IchZiehAus{MyMax: state.LeftMax}, state.Left)
	}
}

func (state *NodeActor) replaceParent(context actor.Context) {
	msg := context.Message().(*messages.IchZiehAus)
	newPID := context.Sender()
	var oldPID *actor.PID
	if msg.MyMax > state.LeftMax {
		// rechts
		oldPID = state.Right
		state.Right = newPID
	} else {
		oldPID = state.Left
		state.Left = newPID
		state.LeftMax = msg.MyMax
	}
	oldPID.Stop()
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
	case *messages.IchZiehAus:
		state.replaceParent(context)
	}
}

func sortKeys(m map[int32]string) []int {
	var keys []int
	for k := range m {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	return keys
}

func sortMap(m map[int32]string) map[int32]string {
	sortedKeys := sortKeys(m)
	mapSorted := make(map[int32]string)

	for i := 0; i < len(m); i++ {
		mapSorted[int32(sortedKeys[i])] = m[int32(sortedKeys[i])]
	}
	return mapSorted
}

func split(m map[int32]string) (leftMap map[int32]string, rightMap map[int32]string, leftMax int) {
	sortedKeys := sortKeys(m)
	lengthMap := len(m) / 2
	leftMap = make(map[int32]string)
	rightMap = make(map[int32]string)

	for i := 0; i <= lengthMap; i++ {
		leftMap[int32(sortedKeys[i])] = m[int32(sortedKeys[i])]
	}
	for i := lengthMap + 1; i < len(m); i++ {
		rightMap[int32(sortedKeys[i])] = m[int32(sortedKeys[i])]
	}
	return leftMap, rightMap, sortedKeys[lengthMap]
}
