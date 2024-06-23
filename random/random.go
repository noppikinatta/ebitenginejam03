package random

import (
	"math/rand/v2"
	"time"
)

var rndCount byte
var chacha8Base [32]byte = [32]byte{
	3, 14, 15, 92, 65, 35, 89, 79,
	32, 38, 64, 26, 43, 38, 32, 79,
	50, 28, 84, 19, 71, 69, 39, 93,
	75, 10, 58, 20, 97, 49, 44, 59,
}

func Source() rand.Source {
	r := byte(time.Now().UnixNano() % 256)
	r += rndCount
	rndCount++
	c8src := [32]byte{}
	for i := range c8src {
		c8src[i] = chacha8Base[i] + r
	}

	return rand.NewChaCha8(c8src)
}
