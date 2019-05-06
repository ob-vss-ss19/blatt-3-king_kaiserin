package main

import (
	"flag"
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
	"sync"
)

type CLINode struct {
	waitgroup *sync.WaitGroup
}

type HelloMsg struct {}

func (state *CLINode) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *HelloMsg:
		fmt.Println("got string %v", msg)
	}
}


func main() {

	flagCreateTree := flag.Bool("newTree", false, "creates new tree, prints out id and token")
	flagLeafSize := flag.Int("size", 1, "size of a leaf")

	//flagID := flag.Int("ID", 1, "ID of the Tree")
	//flagToken := flag.Int("token", 1, "Token of the Tree")

	flagInsert := flag.Bool("insert", false, "insert new value into the tree")
	flagSearch := flag.Bool("search", false, "search value for a key")
	flagDelete := flag.Bool("delete", false, "delete value and key from tree")
	flagTraverse := flag.Bool("traverse", false, "go through tree and get sorted key-value-Pairs")

	flagKey := flag.Int("key", 1, "Key which is needed for Insert/Search/Delete")
	flagValue := flag.String("value", "", "Vale which is needed to insert new key-value-Pair")

	flag.Parse()

	var msg interface{}
	switch  {
	case *flagCreateTree:
		msg = &messages.CheckLeftMax{MaxKey: int32(*flagLeafSize)}
	case *flagTraverse:
		//traverse msg neu mit id und token
		msg = &messages.Traverse{}
	case *flagInsert:
		msg = &messages.Insert{Key: int32(*flagKey), Value: *flagValue }
	case *flagDelete:
		msg = &messages.Delete{Key: int32(*flagKey)}
	case *flagSearch:
		msg = &messages.Search{Key: int32(*flagKey)}

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


