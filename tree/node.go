package tree

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"sort"
	"time"
)

type Insert struct {
	Key   int
	Value string
}

type InsertMap struct {
	inserts map[int]string
}

type Delete struct {
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
}

type NodeActor struct {
	Left      *actor.PID
	LeftMax   int
	Right     *actor.PID
	Leaves    map[int]string
	MaxLeaves int
	Parent    *actor.PID
}

func (state *NodeActor) traverse(context actor.Context){
	if state.Left != nil {
		leftSide, _err := context.RequestFuture(state.Left, &Traverse{}, 1*time.Second).Result()
		rightSide, _err := context.RequestFuture(state.Right, &Traverse{}, 1*time.Second).Result()
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
	msg := context.Message().(*Insert)
	if state.Left != nil {
		if msg.Key > state.LeftMax {
			// msg an rechten Node, dass wert einfuegen
			context.Send(state.Right, &Insert{msg.Key, msg.Value})
		} else {
			// msg an linken teilbaum dass er sich drum kuemmert
			context.Send(state.Left, &Insert{msg.Key, msg.Value})
		}
	} else {
		if state.Leaves == nil {
			state.Leaves = make(map[int]string)
		}
		// pruefen ob map schon voll
		if len(state.Leaves) == state.MaxLeaves {
			// split
			props := actor.PropsFromProducer(func() actor.Actor {
				return &NodeActor{nil, -1, nil, nil, 4, context.Self()}
			})
			state.Left = context.Spawn(props)
			state.Right = context.Spawn(props)
			state.Leaves[msg.Key] = msg.Value
			leftMap, rightMap, leftmaximum := split(state.Leaves)
			state.LeftMax = leftmaximum
			context.Send(state.Left, &InsertMap{leftMap})
			context.Send(state.Right, &InsertMap{rightMap})
			state.Leaves = nil
		} else {
			state.Leaves[msg.Key] = msg.Value
		}
	}
}

func (state *NodeActor) search(context actor.Context) {
	msg := context.Message().(*Search)
	if state.Left != nil {
		if msg.Key < state.LeftMax {
			// an linken weiterschicken
			//context.RequestWithCustomSender()
			context.RequestWithCustomSender(state.Left, &Search{msg.Key}, context.Sender())
		} else {
			// an rechten weiterschicken
			context.RequestWithCustomSender(state.Right, &Search{msg.Key}, context.Sender())
		}
	} else {
		// bei mir oder gar nicht existent
		if value, ok := state.Leaves[msg.Key]; ok {
			context.Respond(&scottyBeamMichHoch{msg.Key, value, ok})
		} else {
			context.Respond(&scottyBeamMichHoch{msg.Key,value, ok})
		}
	}
}

func (state *NodeActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *Insert:
		state.insert(context)
	case *InsertMap:
		state.Leaves = msg.inserts
	case *Delete:
		fmt.Printf("Hello, I will kill you now!")
	case *Search:
		state.search(context)
	case *Traverse:
		state.traverse(context)
	case *scottyBeamMichHoch:
		if msg.ok {
			fmt.Printf("For the key '%v' there is a value '%v'! \n", msg.key, msg.value)
		} else {
			fmt.Printf("For the key '%v' there is NO value! \n", msg.key)
		}

	}
}

func sortMap(m map[int]string) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func split(m map[int]string) (leftMap map[int]string, rightMap map[int]string, leftMax int) {
	sortedKeys := sortMap(m)
	lengthMap := len(m) / 2
	leftMap = make(map[int]string)
	rightMap = make(map[int]string)

	for i := 0; i <= lengthMap; i++ {
		leftMap[sortedKeys[i]] = m[sortedKeys[i]]
	}
	for i := lengthMap + 1; i < len(m); i++ {
		rightMap[sortedKeys[i]] = m[sortedKeys[i]]
	}
	return leftMap, rightMap, sortedKeys[lengthMap]
}
