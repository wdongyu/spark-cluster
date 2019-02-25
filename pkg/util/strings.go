package util

import (
	"math/rand"
	"time"
)

var (
	letterRunes       = []rune("abcdefghijklmnopqrstuvwxyz")
	letterNumberRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterNumberRunes[rand.Intn(len(letterNumberRunes))]
	}
	return string(b)
}
