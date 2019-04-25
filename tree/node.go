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
	keys []int
}

type ShowTree struct {
}

type NodeActor struct {
	Left      *actor.PID
	LeftMax   int
	Right     *actor.PID
	Leaves    map[int]string
	MaxLeaves int
	Parent    *actor.PID
}

func (state *NodeActor) sortedOutput(context actor.Context){
	fmt.Printf("go through left, then right")
	if state.Left != nil {
		leftSide, _err := context.RequestFuture(state.Left, &Traverse{}, 1*time.Second).Result()
		rightSide, _err := context.RequestFuture(state.Right, &Traverse{}, 1*time.Second).Result()
		if _err == nil {
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
			fmt.Printf("All keys in Tree sorted: %v", full)
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
			fmt.Printf("current left baum \n")
		} else {
			// msg an linken teilbaum dass er sich drum kuemmert
			fmt.Printf("muss in linken teilbaum \n")
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
			fmt.Printf("leaves %v", state.Leaves)
			leftMap, rightMap, leftmaximum := split(state.Leaves)
			state.LeftMax = leftmaximum
			context.Send(state.Left, &InsertMap{leftMap})
			context.Send(state.Right, &InsertMap{rightMap})
			state.Leaves = nil
		} else {
			state.Leaves[msg.Key] = msg.Value
		}
	}
	fmt.Printf("Current map leaves: %v \n", state.Leaves)
}

func (state *NodeActor) search(context actor.Context) {
	msg := context.Message().(*Search)
	if state.Left != nil {
		if msg.Key < state.LeftMax {
			// an linken weiterschicken
			//context.RequestWithCustomSender()
			context.Send(state.Left, &Search{msg.Key})
		} else {
			// an rechten weiterschicken
			context.Send(state.Right, &Search{msg.Key})
		}
	} else {
		// bei mir oder gar nicht existent
		if value, ok := state.Leaves[msg.Key]; ok {
			fmt.Printf("tmp found: %v \n", value)
		} else {
			fmt.Printf("not found \n")
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
		state.sortedOutput(context)
	case *scottyBeamMichHoch:
		fmt.Printf("schick die keys zurueck! \n")
	//case *ShowTree:
	//	if state.Left != nil {
	//		context.Send(state.Left, &ShowTree{})
	//		context.Send(state.Right, &ShowTree{})
	//	} else {
	//		fmt.Printf("my map: %v \n", state.Leaves)
	//	}
	}
}

func sortMap(m map[int]string) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	fmt.Printf("sort %v \n", keys)
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
	fmt.Printf("split left %v right %v leftmax %v \n", leftMap, rightMap, sortedKeys[lengthMap])
	return leftMap, rightMap, sortedKeys[lengthMap]
}
