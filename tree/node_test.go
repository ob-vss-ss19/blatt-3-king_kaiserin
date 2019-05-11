package tree

import (
	"fmt"
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
)

func TestCreateEmptyRoot(t *testing.T) {
	fmt.Println("Starting Test")
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext
	fmt.Println("Created Actor")

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		fmt.Println("no Error")
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			fmt.Println("Right Messges Typ")
			resSilce := response.Sorted
			if len(resSilce) != 0 {
				t.Errorf("Expected length of Sorted-Slice was %v but was %v instead. \n", 0, len(resSilce))
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
	fmt.Println("Ending Test")
}

func TestAddValueToRoot(t *testing.T) {
	fmt.Println("Starting Test")
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext
	fmt.Println("Created Actor")

	context.Send(root, &messages.Insert{Key: 2, Value: "zwei"})

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		fmt.Println("no Error")
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			fmt.Println("Right Messges Typ")
			resSilce := response.Sorted
			if len(resSilce) != 1 {
				t.Errorf("Expected length of Sorted-Slice was %v but was %v instead. \n", 1, len(resSilce))
			}
			if resSilce[0].Key != 2 {
				t.Errorf("Expected key was %v but was %v instead. \n", 2, resSilce[0].Key)
			}
			if resSilce[0].Value != "zwei" {
				t.Errorf("Expected value was %v but was %v instead. \n", "zwei", resSilce[0].Value)
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
	fmt.Println("Ending Test")
}

func TestRootSplit(t *testing.T) {
	fmt.Println("Starting Test")
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext
	fmt.Println("Created Actor")

	context.Send(root, &messages.Insert{Key: 2, Value: "zwei"})
	context.Send(root, &messages.Insert{Key: 4, Value: "vier"})

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		fmt.Println("no Error")
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			fmt.Println("Right Messges Typ")
			resSilce := response.Sorted
			if len(resSilce) != 2 {
				t.Errorf("Expected length of Sorted-Slice was %v but was %v instead. \n", 2, len(resSilce))
			}
			if resSilce[0].Key != 2 {
				t.Errorf("Expected key was %v but was %v instead. \n", 2, resSilce[0].Key)
			}
			if resSilce[0].Value != "zwei" {
				t.Errorf("Expected value was %v but was %v instead. \n", "zwei", resSilce[0].Value)
			}
			if resSilce[1].Key != 4 {
				t.Errorf("Expected key was %v but was %v instead. \n", 2, resSilce[0].Key)
			}
			if resSilce[1].Value != "vier" {
				t.Errorf("Expected value was %v but was %v instead. \n", "zwei", resSilce[0].Value)
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
	fmt.Println("Ending Test")
}

func TestLargerTree(t *testing.T) {
	fmt.Println("Starting Test")
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 2, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext
	fmt.Println("Created Actor")

	keys := []int32{int32(5), int32(8), int32(10), int32(20), int32(30), int32(40)}
	values := []string{"fuenf", "acht", "zehn", "zwanzig", "dreissig", "vierzieg"}

	context.Send(root, &messages.Insert{Key: keys[2], Value: values[2]})
	context.Send(root, &messages.Insert{Key: keys[3], Value: values[3]})
	context.Send(root, &messages.Insert{Key: keys[1], Value: values[1]})
	context.Send(root, &messages.Insert{Key: keys[0], Value: values[0]})
	context.Send(root, &messages.Insert{Key: keys[4], Value: values[4]})
	context.Send(root, &messages.Insert{Key: keys[5], Value: values[5]})

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		fmt.Println("no Error")
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			fmt.Println("Right Messges Typ")
			resSilce := response.Sorted
			if len(resSilce) != len(keys) {
				t.Errorf("Expected length of Sorted-Slice was %v but was %v instead. \n", len(keys), len(resSilce))
			}
			for i := range keys {
				if resSilce[i].Key != keys[i] {
					t.Errorf("Expected key was %v but was %v instead. \n", keys[i], resSilce[i].Key)
				}
				if resSilce[i].Value != values[i] {
					t.Errorf("Expected value was %v but was %v instead. \n", values[i], resSilce[i].Value)
				}
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
	fmt.Println("Ending Test")
}

func TestSearchLeft(t *testing.T) {
	fmt.Println("Starting Test")
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext
	fmt.Println("Created Actor")

	context.Send(root, &messages.Insert{Key: 2, Value: "zwei"})
	context.Send(root, &messages.Insert{Key: 4, Value: "vier"})

	future := context.RequestFuture(root, &messages.Search{Key: 2}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		fmt.Println("no Error")
		response, ok := res.(*messages.ScottyBeamMichHoch)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			fmt.Println("Right Messges Typ")
			if !response.Ok {
				t.Errorf("Expected to find value but didn't. \n")
			}
			if response.Key != 2 {
				t.Errorf("Expected key was %v but was %v instead. \n", 2, response.Key)
			}
			if response.Value != "zwei" {
				t.Errorf("Expected value was %v but was %v instead. \n", "zwei", response.Value)
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
	fmt.Println("Ending Test")
}

func TestSearchRight(t *testing.T) {
	fmt.Println("Starting Test")
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext
	fmt.Println("Created Actor")

	context.Send(root, &messages.Insert{Key: 2, Value: "zwei"})
	context.Send(root, &messages.Insert{Key: 4, Value: "vier"})

	future := context.RequestFuture(root, &messages.Search{Key: 4}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		fmt.Println("no Error")
		response, ok := res.(*messages.ScottyBeamMichHoch)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			fmt.Println("Right Messges Typ")
			if !response.Ok {
				t.Errorf("Expected to find value but didn't. \n")
			}
			if response.Key != 4 {
				t.Errorf("Expected key was %v but was %v instead. \n", 4, response.Key)
			}
			if response.Value != "vier" {
				t.Errorf("Expected value was %v but was %v instead. \n", "vier", response.Value)
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
	fmt.Println("Ending Test")
}

func TestSearchFail(t *testing.T) {
	fmt.Println("Starting Test")
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext
	fmt.Println("Created Actor")

	future := context.RequestFuture(root, &messages.Search{Key: 4}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		fmt.Println("no Error")
		response, ok := res.(*messages.ScottyBeamMichHoch)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			fmt.Println("Right Messges Typ")
			if response.Ok {
				t.Errorf("Didn't expected to find something but found key: %v and value: %v  \n", response.Key, response.Value)
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
	fmt.Println("Ending Test")
}

func TestDeleteOneLeave(t *testing.T) {
	fmt.Println("Starting Test")
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext
	fmt.Println("Created Actor")

	context.Send(root, &messages.Insert{Key: 2, Value: "zwei"})
	context.Send(root, &messages.Delete{Key: 2})

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		fmt.Println("no Error")
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			fmt.Println("Right Messges Typ")
			resSilce := response.Sorted
			if len(resSilce) != 0 {
				t.Errorf("Expected no elemnts but found: %v \n", resSilce)
			}
		}
	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
	fmt.Println("Ending Test")
}

func TestDeleteLargerLeave(t *testing.T) {
	fmt.Println("Starting Test")
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 2, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext
	fmt.Println("Created Actor")

	keys := []int32{int32(5), int32(9), int32(10)}
	values := []string{"fuenf", "neun", "zehn"}

	context.Send(root, &messages.Insert{Key: keys[0], Value: values[0]})
	context.Send(root, &messages.Insert{Key: keys[1], Value: values[1]})
	context.Send(root, &messages.Insert{Key: keys[2], Value: values[2]})
	context.Send(root, &messages.Insert{Key: 1, Value: "eins"})

	context.Send(root, &messages.Delete{Key: 1})

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		fmt.Println("no Error")
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			fmt.Println("Right Messges Typ")
			resSilce := response.Sorted
			if len(resSilce) != len(keys) {
				t.Errorf("Expected length of Sorted-Slice was %v but was %v instead. \n", len(keys), len(resSilce))
			}
			for i := range keys {
				if resSilce[i].Key != keys[i] {
					t.Errorf("Expected key was %v but was %v instead. \n", keys[i], resSilce[i].Key)
				}
				if resSilce[i].Value != values[i] {
					t.Errorf("Expected value was %v but was %v instead. \n", values[i], resSilce[i].Value)
				}
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
	fmt.Println("Ending Test")
}

func TestDeleteRightChild(t *testing.T) {
	fmt.Println("Starting Test")
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext
	fmt.Println("Created Actor")

	keys := []int32{int32(1), int32(2), int32(8), int32(9)}
	values := []string{"eins", "zwei", "acht", "neun"}

	context.Send(root, &messages.Insert{Key: keys[1], Value: values[1]})
	context.Send(root, &messages.Insert{Key: 10, Value: "zehn"})
	context.Send(root, &messages.Insert{Key: keys[0], Value: values[0]})
	context.Send(root, &messages.Insert{Key: keys[3], Value: values[3]})
	context.Send(root, &messages.Insert{Key: keys[2], Value: values[2]})

	context.Send(root, &messages.Delete{Key: 10})

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		fmt.Println("no Error")
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			fmt.Println("Right Messges Typ")
			resSilce := response.Sorted
			if len(resSilce) != len(keys) {
				t.Errorf("Expected length of Sorted-Slice was %v but was %v instead. \n", len(keys), len(resSilce))
			}
			for i := range keys {
				if resSilce[i].Key != keys[i] {
					t.Errorf("Expected key was %v but was %v instead. \n", keys[i], resSilce[i].Key)
				}
				if resSilce[i].Value != values[i] {
					t.Errorf("Expected value was %v but was %v instead. \n", values[i], resSilce[i].Value)
				}
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
	fmt.Println("Ending Test")
}

func TestDeleteLeftChild(t *testing.T) {
	fmt.Println("Starting Test")
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext
	fmt.Println("Created Actor")

	keys := []int32{int32(9), int32(10)}
	values := []string{"neun", "zehn"}

	context.Send(root, &messages.Insert{Key: keys[1], Value: values[1]})
	context.Send(root, &messages.Insert{Key: keys[0], Value: values[0]})
	context.Send(root, &messages.Insert{Key: 5, Value: "fuenf"})

	context.Send(root, &messages.Delete{Key: 5})

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		fmt.Println("no Error")
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			fmt.Println("Right Messges Typ")
			resSilce := response.Sorted
			fmt.Printf("value we got: %v \n", resSilce)
			if len(resSilce) != len(keys) {
				t.Errorf("Expected length of Sorted-Slice was %v but was %v instead. \n", len(keys), len(resSilce))
			}
			for i := range keys {
				if resSilce[i].Key != keys[i] {
					t.Errorf("Expected key was %v but was %v instead. \n", keys[i], resSilce[i].Key)
				}
				if resSilce[i].Value != values[i] {
					t.Errorf("Expected value was %v but was %v instead. \n", values[i], resSilce[i].Value)
				}
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
	fmt.Println("Ending Test")
}
