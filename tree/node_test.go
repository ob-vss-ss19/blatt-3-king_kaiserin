package tree

import (
	"testing"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/ob-vss-ss19/blatt-3-king_kaiserin/messages"
)

const valueString2 = "zwei"

func TestCreateEmptyRoot(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			resSilce := response.Sorted
			if len(resSilce) != 0 {
				t.Errorf("Expected length of Sorted-Slice was %v but was %v instead. \n", 0, len(resSilce))
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
}

func TestAddValueToRoot(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext

	_ = context.RequestFuture(root, &messages.Insert{Key: 2, Value: valueString2}, 1*time.Second)

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			resSilce := response.Sorted
			if len(resSilce) != 1 {
				t.Errorf("Expected length of Sorted-Slice was %v but was %v instead. \n", 1, len(resSilce))
			}
			if resSilce[0].Key != 2 {
				t.Errorf("Expected key was %v but was %v instead. \n", 2, resSilce[0].Key)
			}
			if resSilce[0].Value != valueString2 {
				t.Errorf("Expected value was %v but was %v instead. \n", valueString2, resSilce[0].Value)
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
}

func TestRootSplit(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext

	_ = context.RequestFuture(root, &messages.Insert{Key: 2, Value: "valueString2"}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: 4, Value: "vier"}, 1*time.Second)

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			resSilce := response.Sorted
			if len(resSilce) != 2 {
				t.Errorf("Expected length of Sorted-Slice was %v but was %v instead. \n", 2, len(resSilce))
			}
			if resSilce[0].Key != 2 {
				t.Errorf("Expected key was %v but was %v instead. \n", 2, resSilce[0].Key)
			}
			if resSilce[0].Value != "valueString2" {
				t.Errorf("Expected value was %v but was %v instead. \n", "valueString2", resSilce[0].Value)
			}
			if resSilce[1].Key != 4 {
				t.Errorf("Expected key was %v but was %v instead. \n", 2, resSilce[0].Key)
			}
			if resSilce[1].Value != "vier" {
				t.Errorf("Expected value was %v but was %v instead. \n", "valueString2", resSilce[0].Value)
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
}

func TestLargerTree(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 2, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext

	keys := []int32{int32(5), int32(8), int32(10), int32(20), int32(30), int32(40)}
	values := []string{"fuenf", "acht", "zehn", "zwanzig", "dreissig", "vierzieg"}

	_ = context.RequestFuture(root, &messages.Insert{Key: keys[2], Value: values[2]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: keys[3], Value: values[3]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: keys[1], Value: values[1]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: keys[0], Value: values[0]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: keys[4], Value: values[4]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: keys[5], Value: values[5]}, 1*time.Second)

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
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
}

func TestSearchLeft(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext

	_ = context.RequestFuture(root, &messages.Insert{Key: 2, Value: "valueString2"}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: 4, Value: "vier"}, 1*time.Second)

	future := context.RequestFuture(root, &messages.Search{Key: 2}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		response, ok := res.(*messages.ScottyBeamMichHoch)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			if !response.Ok {
				t.Errorf("Expected to find value but didn't. \n")
			}
			if response.Key != 2 {
				t.Errorf("Expected key was %v but was %v instead. \n", 2, response.Key)
			}
			if response.Value != "valueString2" {
				t.Errorf("Expected value was %v but was %v instead. \n", "valueString2", response.Value)
			}
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
	future = context.RequestFuture(root, &messages.Search{Key: 4}, 1*time.Second)
	res, err = future.Result()
	if err == nil {
		response, ok := res.(*messages.ScottyBeamMichHoch)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
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
}

func TestSearchFail(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext

	future := context.RequestFuture(root, &messages.Search{Key: 4}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		response, ok := res.(*messages.ScottyBeamMichHoch)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else if response.Ok {
			t.Errorf("Didn't expected to find something but found key: %v and value: %v  \n", response.Key, response.Value)
		}

	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
}

func TestDeleteOneLeave(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext

	_ = context.RequestFuture(root, &messages.Insert{Key: 2, Value: "valueString2"}, 1*time.Second)
	futureDelete := context.RequestFuture(root, &messages.Delete{Key: 2}, 1*time.Second)
	resDelete, errDelete := futureDelete.Result()
	if errDelete == nil {
		_, ok := resDelete.(*messages.BaumFaellt)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		}
	} else {
		t.Errorf("Error getting Future: %v \n", errDelete)
	}

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
			resSilce := response.Sorted
			if len(resSilce) != 0 {
				t.Errorf("Expected no elemnts but found: %v \n", resSilce)
			}
		}
	} else {
		t.Errorf("Error getting Future: %v \n", err)
	}
}

func TestDeleteLargerLeave(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 2, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext

	keys := []int32{int32(5), int32(9), int32(10)}
	values := []string{"fuenf", "neun", "zehn"}

	_ = context.RequestFuture(root, &messages.Insert{Key: keys[0], Value: values[0]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: keys[1], Value: values[1]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: keys[2], Value: values[2]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: 1, Value: "eins"}, 1*time.Second)

	futureDelete := context.RequestFuture(root, &messages.Delete{Key: 1}, 1*time.Second)

	resDelete, errDelete := futureDelete.Result()
	if errDelete == nil {
		resDelete, ok := resDelete.(*messages.DeleteResult)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else if !resDelete.Successful {
			t.Error("Expected successful delete! \n")
		}
	} else {
		t.Errorf("Error getting Future: %v \n", errDelete)
	}

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
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
}

func TestDeleteRightChild(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext

	keys := []int32{int32(1), int32(2), int32(8), int32(9)}
	values := []string{"eins", "valueString2", "acht", "neun"}

	_ = context.RequestFuture(root, &messages.Insert{Key: keys[1], Value: values[1]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: 10, Value: "zehn"}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: keys[0], Value: values[0]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: keys[3], Value: values[3]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: keys[2], Value: values[2]}, 1*time.Second)

	futureDelete := context.RequestFuture(root, &messages.Delete{Key: 10}, 1*time.Second)

	resDelete, errDelete := futureDelete.Result()
	if errDelete == nil {
		resDelete, ok := resDelete.(*messages.DeleteResult)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else if !resDelete.Successful {
			t.Error("Expected successful delete! \n")
		}
	} else {
		t.Errorf("Error getting Future: %v \n", errDelete)
	}

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
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
}

func TestDeleteLeftChild(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext

	keys := []int32{int32(9), int32(10)}
	values := []string{"neun", "zehn"}

	_ = context.RequestFuture(root, &messages.Insert{Key: keys[1], Value: values[1]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: keys[0], Value: values[0]}, 1*time.Second)
	_ = context.RequestFuture(root, &messages.Insert{Key: 5, Value: "fuenf"}, 1*time.Second)

	futureDelete := context.RequestFuture(root, &messages.Delete{Key: 5}, 1*time.Second)

	resDelete, errDelete := futureDelete.Result()
	if errDelete == nil {
		resDelete, ok := resDelete.(*messages.DeleteResult)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else if !resDelete.Successful {
			t.Error("Expected successful delete! \n")
		}
	} else {
		t.Errorf("Error getting Future: %v \n", errDelete)
	}

	future := context.RequestFuture(root, &messages.Traverse{}, 1*time.Second)
	res, err := future.Result()
	if err == nil {
		response, ok := res.(*messages.TraverseResponse)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else {
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
}

func TestDeleteNotExistingKey(t *testing.T) {
	context := actor.EmptyRootContext
	props := actor.PropsFromProducer(func() actor.Actor {
		return &NodeActor{nil, nil, nil, nil, -1, 1, -1, ""}
	})
	root := context.Spawn(props)
	context = actor.EmptyRootContext

	_ = context.RequestFuture(root, &messages.Insert{Key: 2, Value: "valueString2"}, 1*time.Second)
	futureDelete := context.RequestFuture(root, &messages.Delete{Key: 5}, 1*time.Second)
	resDelete, errDelete := futureDelete.Result()
	if errDelete == nil {
		res, ok := resDelete.(*messages.DeleteResult)
		if !ok {
			t.Error("Expected other Msg Type! \n")
		} else if res.Successful {
			t.Error("Expected failing delete!\n")
		}
	} else {
		t.Errorf("Error getting Future: %v \n", errDelete)
	}
}
