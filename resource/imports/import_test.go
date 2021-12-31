package main

import (
	"fmt"
	"testing"
)

func Test_createGroup(t *testing.T) {
	err := createGroup("dx", "test-dx", []string{"10.1.87.56:432"}, 80)
	fmt.Println(err)
}

func Test_createIdc(t *testing.T) {
	err := createIdc(&Idc{
		Name: "test",
		Desc: "测试",
	})
	fmt.Println(err)
}

func Test_parsePools(t *testing.T) {
	err := parsePools()
	fmt.Println(err)
}

func TestPOol(t *testing.T) {
	parsePoolsGroup()

}



//

func Test_readHostGroupRules(t *testing.T) {
	ss ,rr, pp := readHostGroupRules()
	for _, s := range ss {
		fmt.Println(*s)
	}
	fmt.Println(len(ss),"--------")

	for _, r := range rr {
		fmt.Println(*r)
	}
	fmt.Println(len(rr),"--------")
	for k, p := range pp {
		fmt.Println(k,p)
	}
}
