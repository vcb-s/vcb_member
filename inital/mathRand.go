package inital

import (
	mathRand "math/rand"
	"time"
)

func setupMathSeed() {
	mathRand.Seed(time.Now().Unix())
}
