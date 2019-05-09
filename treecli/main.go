package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
)

type CLINode struct {
	waitgroup *sync.WaitGroup
}

func (state *CLINode) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *messages.ScottyBeamMichHoch:
		if msg.Ok {
			fmt.Printf("For the key '%v' there is a value '%v'! \n", msg.Key, msg.Value)
		} else {
			fmt.Printf("For the key '%v' there is NO value! \n", msg.Key)
		}
		state.waitgroup.Done()
	case *messages.TraverseResponse:
		fmt.Printf("All keys in Tree sorted: %v\n", msg.Sorted)
		state.waitgroup.Done()
	}
}

func main() {

	flagCreateTree := flag.Bool("newTree", false, "creates new tree, prints out id and token")
	flagLeafSize := flag.Int("size", 1, "size of a leaf")

	flagID := flag.Int("ID", 1, "ID of the Tree")
	flagToken := flag.String("token", "", "Token of the Tree")

	flagInsert := flag.Bool("insert", false, "insert new value into the tree")
	flagSearch := flag.Bool("search", false, "search value for a key")
	flagDelete := flag.Bool("delete", false, "delete value and key from tree")
	flagTraverse := flag.Bool("traverse", false, "go through tree and get sorted key-value-Pairs")

	flagKey := flag.Int("key", 1, "Key which is needed for Insert/Search/Delete")
	flagValue := flag.String("value", "", "Vale which is needed to insert new key-value-Pair")

	flag.Parse()

	var msg interface{}
	switch {
	case *flagCreateTree:
		msg = &messages.PflanzBaum{MaxLeaves: int32(*flagLeafSize)}
	case *flagTraverse:
		find := &messages.Tree{ID: int32(*flagID), Token: *flagToken}
		msg = &messages.TraverseCLI{Find: find}
	case *flagInsert:
		find := &messages.Tree{ID: int32(*flagID), Token: *flagToken}
		msg = &messages.InsertCLI{Find: find, Key: int32(*flagKey), Value: *flagValue}
	case *flagDelete:
		find := &messages.Tree{ID: int32(*flagID), Token: *flagToken}
		msg = &messages.DeleteCLI{Find: find, Key: int32(*flagKey)}
	case *flagSearch:
		find := &messages.Tree{ID: int32(*flagID), Token: *flagToken}
		msg = &messages.SearchCLI{Find: find, Key: int32(*flagKey)}
	}

	remote.Start("localhost:8091")
	var waitgroup sync.WaitGroup

	props := actor.PropsFromProducer(
		func() actor.Actor {
			waitgroup.Add(1)
			return &CLINode{&waitgroup}
		})

	cli := actor.Spawn(props)
	context := actor.EmptyRootContext
	remote := actor.NewPID("localhost:8090", "service")

	//msg := messages.CheckLeftMax{MaxKey: 5}
	fmt.Printf("kurz vor message \n")
	context.RequestWithCustomSender(remote, msg, cli)
	fmt.Printf("message gesendet \n")

	waitgroup.Wait()
}
