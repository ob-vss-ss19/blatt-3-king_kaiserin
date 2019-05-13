package main

//import (
//	"github.com/AsynkronIT/goconsole"
//	"github.com/AsynkronIT/protoactor-go/actor"
//	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
//	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/tree"
//)
//
//func main() {
//	context := actor.EmptyRootContext
//	props := actor.PropsFromProducer(func() actor.Actor {
//		return &tree.NodeActor{nil, nil, nil, nil, -1, 2, -1, ""}
//	})
//	pid := context.Spawn(props)
//	context.Send(pid, &messages.Insert{Key: 5, Value: "five"})
//	context.Send(pid, &messages.Insert{Key: 7, Value: "seven"})
//	context.Send(pid, &messages.Insert{Key: 9, Value: "nine"})
//
//	context.Send(pid, &messages.Traverse{})
//	context.RequestWithCustomSender(pid, &messages.Delete{Key: 9}, pid)
//	context.Send(pid, &messages.Traverse{})
//
//	console.ReadLine()
//}
