package tools

import (
	"math/rand"
)

type Swap func(i,j int)
func Shuffle(swap Swap,length int ) {
	for i := 0; i < length; i++ {
		a := rand.Int() % length
		swap(i,a)
	}
}

var(
	seed  int64 = 0
)

func shuffle2(arr []int){
	length:= len(arr)
	for i := 0; i < length; i++ {
		a := rand.Int() % length
		arr[i],arr[a] = arr[a],arr[i]
	}
}
