package randext

import (
	"math/rand"
)

func Bool() bool {
	return rand.Intn(2) == 0
}
