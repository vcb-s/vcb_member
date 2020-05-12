package inital

import (
	mathRand "math/rand"
	"time"
)

func init() {
	mathRand.Seed(time.Now().Unix())
}
