package core

import (
	"fmt"
	"testing"
)

func Test_newWeightLoopBalancer(t *testing.T) {
	lw := newWeightLoopBalancer([]*Target{
		{Addr:   "127.0.0.1", Weight: 100,},
		{Addr:   "127.0.0.2", Weight: 100,},
		{Addr:   "127.0.0.3", Weight: 100,},
		{Addr:   "127.0.0.4", Weight: 50,},
	},4)

	fmt.Println(lw.getAddress())
	fmt.Println(lw.getAddress())
	fmt.Println(lw.getAddress())
	fmt.Println(lw.getAddress())
	fmt.Println(lw.getAddress())
	fmt.Println(lw.getAddress())
	fmt.Println(lw.getAddress())
	fmt.Println(lw.getAddress())
	count := map[string]int{}
	for i:=0;i<10000;i++{
		for _, t := range lw.getAddress() {
			count[t.Addr]++
		}
	}
	fmt.Println(count)
}

func BenchmarkLoopBanalcer(b *testing.B) {
	lw := newWeightLoopBalancer([]*Target{
		{Addr:   "127.0.0.1", Weight: 100,},
		{Addr:   "127.0.0.2", Weight: 100,},
		{Addr:   "127.0.0.3", Weight: 100,},
		{Addr:   "127.0.0.4", Weight: 100,},
	},4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next(){
			lw.getAddress()
		}
	})
}
