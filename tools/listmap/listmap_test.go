package listmap

import (
	"fmt"
	"strconv"
	"testing"
)

func TestListMap_Set(t *testing.T) {
	m:=&ListMap{}
	for i:=0;i<10000000;i++{
		m.Set(strconv.Itoa(i %10),i)
	}
}

func TestListMap_Get(t *testing.T) {
	m:=&ListMap{}
	for i:=0;i<10;i++{
		m.Set(strconv.Itoa(i %5),i)
	}
	var v interface{}
	for i:=0;i<10000000;i++{
		v = m.Get(strconv.Itoa(i %5))
	}
	fmt.Println(v)
}

func TestListMap_Set2(t *testing.T) {
	m:=make(map[string]interface{})
	for i:=0;i<10;i++{
		m[strconv.Itoa(i%20)] = i
	}
	var v interface{}
	for i:=0;i<10000000;i++{
		v = m[strconv.Itoa(i%20)]
	}
	fmt.Println(v)
}

