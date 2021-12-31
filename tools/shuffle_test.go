package tools

import (
	"fmt"
	"math/rand"
	"testing"
	"unsafe"
)

func TestShuffle(t *testing.T) {
	c:=5
	for i:=0;i<5;i++{
		c:=rand.Int()
		fmt.Println(c,uintptr(unsafe.Pointer(&c)),uintptr(unsafe.Pointer(&i)))
	}
	fmt.Println(c,uintptr(unsafe.Pointer(&c)))
}

var arr = []int{1,2,3,4,5}

func BenchmarkShuffle(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Shuffle(func(i, j int) {
			arr[i],arr[j] = arr[j],arr[i]
		}, len(arr))
	}
}


func BenchmarkShuffle2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		shuffle2(arr)
	}

}

func TestSplitIp(t *testing.T) {
	fmt.Println(SplitIp("12.21.12.2"))
	fmt.Println(SplitIp("12.21.12.2:98"))
	fmt.Println(SplitIp("12.21.12.2:98r"))
	fmt.Println(SplitIp("[12.21:12:2:]:98"))
	fmt.Println(SplitIp("[12.21:12:2:]"))
	fmt.Println(IsValidIpAddress("1.2.3.4"))
	fmt.Println(IsValidIpAddress("1.2.3.4:90"))
	fmt.Println(IsValidIpAddress("[2400:da00::dbf:0:100]"))
	fmt.Println(IsValidIpAddress("[2400:da00::dbf:0:100]:8000"))
	fmt.Println(IsValidIpAddress("[2400:da00::dbf:0:1000]:8000"))
}



