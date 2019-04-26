package tree

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
	"sort"
	"time"
)

/*type Insert struct {
	Key   int
	Value string
}*/

type InsertMap struct {
	inserts map[int32]string
}

/*type Delete struct {
	Key   int
	Value string
}

type Search struct {
	Key int
}

type Traverse struct{}

type scottyBeamMichHoch struct {
	key int
	value string
	ok bool
}*/

type NodeActor struct {
	Left      *actor.PID
	LeftMax   int32
	Right     *actor.PID
	Leaves    map[int32]string
	MaxLeaves int32
	Parent    *actor.PID
}

func (state *NodeActor) traverse(context actor.Context){
	if state.Left != nil {
		leftSide, _err := context.RequestFuture(state.Left, &messages.Traverse{}, 1*time.Second).Result()
		rightSide, _err := context.RequestFuture(state.Right, &messages.Traverse{}, 1*time.Second).Result()
		if _err != nil {
			println("Error with Futures happened!")
		}

		lSide := leftSide.(*[]int)
		rSide := rightSide.(*[]int)
		var full []int
		full = append(full, *lSide...)
		full = append(full, *rSide...)

		if state.Parent != nil {
			context.Respond(full)
		} else {
			fmt.Printf("All keys in Tree sorted: %v\n", full)
		}
	} else {
		leaves := sortMap(state.Leaves)
		context.Respond(&leaves)
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
		if msg.Key < state.LeftMax {
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
			context.Respond(&messages.ScottyBeamMichHoch{Key: msg.Key,Value: value, Ok: ok})
		}
	}
}

func (state *NodeActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.Insert:
		state.insert(context)
	case *InsertMap:
		state.Leaves = msg.inserts
	case *messages.Delete:
		fmt.Printf("Hello, I will kill you now!")
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

	}
}

func sortMap(m map[int32]string) []int {
	var keys []int
	for k := range m {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	return keys
}

func split(m map[int32]string) (leftMap map[int32]string, rightMap map[int32]string, leftMax int) {
	sortedKeys := sortMap(m)
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
