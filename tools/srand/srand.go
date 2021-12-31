package srand

import "sync/atomic"

var seed int64 = 1

var (
	randList = []int{}
	randMax  = int64(10000)
)

func init() {
	randList = make([]int, randMax)
	for i := 0; i < len(randList); i++ {
		randList[i] = i
	}
}

func Rand(n int64) int64 {
	if n >= randMax {
		n = randMax
	}
	return atomic.AddInt64(&seed, 1) % n
}
