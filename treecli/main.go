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
		//state.waitgroup.Done()
	case *messages.TraverseResponse:
		fmt.Printf("All keys in Tree sorted: %v\n", msg.Sorted)
		//state.waitgroup.Done()
	case *messages.BaumFaellt:
		fmt.Printf("loesche tree mit id %v und token %v", msg.ID, msg.Token)
		remote := actor.NewPID("localhost:8090", "service")
		context.Send(remote, &messages.BaumFaellt{ID: msg.ID, Token: msg.Token})
		//state.waitgroup.Done()
	case *messages.PflanzBaumResponse:
		fmt.Printf("Created a new Tree with ID: %v and Token: %v", msg.ID, msg.Token)
		//state.waitgroup.Done()
	case *messages.DeleteResult:
		if msg.Successful {
			fmt.Printf("deleting was successful! \n")
		} else {
			fmt.Printf("deleting was NOT successful! The given key does not exist.\n")
		}
		//state.waitgroup.Done()
	case *messages.TreeNotFound:
		fmt.Printf("Tree with token %v and pid %v not found!\n", msg.NotFound.Token, msg.NotFound.ID)
		//state.waitgroup.Done()
	case *messages.InsertResult:
		if msg.Successful {
			fmt.Printf("inserting was successful! \n")
		} else {
			fmt.Printf("inserting was NOT successful!\n")
		}
		//state.waitgroup.Done()
	case *messages.DeleteTreeRespond:
		if msg.Delete {
			fmt.Printf("Tree was felled!\n")
		} else {
			fmt.Printf("If you really want to delete the tree, send this command once more.\n")
		}
		//state.waitgroup.Done()
	}
}

func main() {

	flagBind := flag.String("bind", "localhost:8091", "Adresse to bind CLI")
	flagRemote := flag.String("remote", "localhost:8090", "Adresse to bind Service")

	flagNameCli := flag.String("nameCLI", "treecli", "Name for the CLI")
	flagNameService := flag.String("nameService", "treeservice", "Name for the Service")

	flagCreateTree := flag.Bool("newTree", false, "creates new tree, prints out id and token")
	flagLeafSize := flag.Int("size", 1, "size of a leaf")

	flagID := flag.Int("ID", 1, "ID of the Tree")
	flagToken := flag.String("token", "", "Token of the Tree")

	flagInsert := flag.Bool("insert", false, "insert new value into the tree")
	flagSearch := flag.Bool("search", false, "search value for a key")
	flagDelete := flag.Bool("delete", false, "delete value and key from tree")
	flagTraverse := flag.Bool("traverse", false, "go through tree and get sorted key-value-Pairs")
	flagDeleteTree := flag.Bool("deleteTree", false, "delete whole Tree")

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
	case *flagDeleteTree:
		find := &messages.Tree{ID: int32(*flagID), Token: *flagToken}
		msg = &messages.DeleteTree{Delete: find}
	}
	if msg != nil {
		remote.Start(*flagBind)
		var waitgroup sync.WaitGroup

		props := actor.PropsFromProducer(
			func() actor.Actor {
				waitgroup.Add(1)
				return &CLINode{&waitgroup}
			})

		cli := actor.Spawn(props)
		remote.Register(*flagNameCli, props)
		context := actor.EmptyRootContext
		remote := actor.NewPID(*flagRemote, *flagNameService)

		context.RequestWithCustomSender(remote, msg, cli)

		waitgroup.Wait()
	} else {
		flag.Usage()
	}

}
