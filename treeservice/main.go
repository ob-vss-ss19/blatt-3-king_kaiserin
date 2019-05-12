package main

import (
	"fmt"
	"sync"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"

	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/treeservice/service"
)

func main() {
	fmt.Printf("Hello Tree-Service!!\n\n")

	remote.Start("localhost:8090")
	var waitgroup sync.WaitGroup

	props := actor.PropsFromProducer(
		func() actor.Actor {
			waitgroup.Add(1)
			return &service.NodeService{Waitgroup: &waitgroup, Roots: make(map[int32]*service.Validation), MarkedForDelete: make(map[int32]*service.Validation), NextID: 1001}
		})

	pid, err := actor.SpawnNamed(props, "service")
	if err == nil {
		fmt.Printf("started %v", *pid)
		waitgroup.Wait()
	} else {
		fmt.Printf("error %v", err.Error())
	}
}
