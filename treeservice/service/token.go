package service

import (
	"math/rand"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyz" + "0123456789"

func CreateToken(length int) string {
	var s = rand.NewSource(time.Now().Unix())
	var r = rand.New(s)

	tokenslice := make([]byte, length)
	for i := range tokenslice {
		tokenslice[i] = chars[r.Intn(len(chars))]
	}
	return string(tokenslice)
}
