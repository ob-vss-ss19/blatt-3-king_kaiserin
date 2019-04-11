package main

import(
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/goconsole"
)


type hello struct{ Who string }

type helloActor struct{}

func (state *helloActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
case *hello:
fmt.Printf("Hello %v\n", msg.Who)
}
}
func main() {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &helloActor{}
	})
	pid := context.Spawn(props)
	context.Send(pid, &hello{Who: "Roger"})
	console.ReadLine()
}
