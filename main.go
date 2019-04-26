package main

import (
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/tree"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
)

func main() {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &tree.NodeActor{nil, -1, nil, nil, 4, nil}
	})
	pid := context.Spawn(props)
	context.Send(pid, &messages.Insert{Key: 5, Value: "five"})
	context.Send(pid, &messages.Insert{Key: 7, Value: "seven"})
	context.Send(pid, &messages.Insert{Key: 9, Value: "nine"})
	context.Send(pid, &messages.Insert{Key: 4, Value: "four"})
	context.Send(pid, &messages.Insert{Key: 6, Value: "six"})
	context.Send(pid, &messages.Insert{Key: 8, Value: "eight"})
	context.RequestWithCustomSender(pid, &messages.Search{Key: 5}, pid)
	context.RequestWithCustomSender(pid, &messages.Search{Key: 3}, pid)

	context.Send(pid, &messages.Traverse{})

	console.ReadLine()
}
