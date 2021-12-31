package dao

import (
	"database/sql"
	"fmt"
	"testing"
	"unsafe"
)

func Test_getMaxEventId(t *testing.T) {
	var err error
	testdb, err = sql.Open("postgres", "host=10.1.87.70 port=55432 dbname=nameserver user=nameserver sslmode=disable")
	if err != nil {
		panic(err)
	}
	fmt.Println(GetMaxEventId(testdb))
}

func eqString(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	ptr1 := (*[]int)(unsafe.Pointer(&s1))
	ptr2 := (*[]int)(unsafe.Pointer(&s2))
	length := (len(s1) +7)/8
	for i := 0; i < length ; i++ {
		if (*ptr1)[i] != (*ptr2)[i]{
			return false
		}
	}
	return true
}

func eqStrint2(s1,s2 string)bool{
	return s1 == s2
}

func TestString(t *testing.T){
	fmt.Println(eqString("22222222222222222222","22222222222222222222"))
}

func BenchmarkClusterEventDao_GetMaxEventId(b *testing.B) {
	for i:=0;i<b.N;i++{
		eqString("22222222222222222222","22222222222222222222")
	}
}

func BenchmarkClusterEventDao_GetMaxEventId2(b *testing.B) {
	for i:=0;i<b.N;i++{
		eqStrint2("22222222222222222222","22222222222222222222")
	}
}


