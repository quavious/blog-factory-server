package utils

import (
	"math/rand"
	"time"
)

const charSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const length = len(charSet)

func NewEmailToken() string {
	rand.Seed(time.Now().Unix())
	token := ""
	for i := 0; i < 6; i++ {
		index := rand.Intn(length)
		v := charSet[index]
		token += string(v)
	}
	return token
}
