package main

import (
	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"./tree"
)

func main() {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &tree.NodeActor{nil, -1, nil, nil, 4}
	})
	pid := context.Spawn(props)
	context.Send(pid, &tree.Insert{5, "five"})
	context.Send(pid, &tree.Insert{7, "seven"})
	context.Send(pid, &tree.Insert{9, "nine"})
	context.Send(pid, &tree.Insert{4, "four"})
	context.Send(pid, &tree.Insert{6, "six"})
	context.Send(pid, &tree.Insert{8, "eight"})
	context.Send(pid, &tree.ShowTree{})

	console.ReadLine()
}
