package hash

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestStringI64(t *testing.T) {

}

func BenchmarkStringI64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringI64("hello world")
	}
}

func BenchmarkHashMap_Set(b *testing.B) {
	m := newHashMap(16)
	for i := 0; i < b.N; i++ {
		m.Set("12", 1)
	}
}

func BenchmarkHashMap_Set2(b *testing.B) {
	m := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		m["12"] = 1
	}

	ml := sync.Mutex{}
	ml.Lock()

}

func TestHashMap_Get(t *testing.T) {
	laddrm, err := net.ResolveTCPAddr("tcp", ":")
	if err != nil {
		panic(err)
	}

	remoteAddr, err := net.ResolveTCPAddr("tcp", "10.1.87.70:8888")
	if err != nil {
		panic(err)
	}
	fmt.Println(err, laddrm, remoteAddr)

	conn,err:=net.DialTCP("tcp",laddrm,remoteAddr)
	if err != nil{
		panic(err)
	}
	fmt.Println(conn.LocalAddr().String())
	conn.Write([]byte("hello"))
}


