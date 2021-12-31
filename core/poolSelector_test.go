package core

import (
	"fmt"
	"testing"
)

func Test_weightPoolSelector_selectAddrs(t *testing.T) {
	ps := weightPoolSelector{

	}
	p:=&pool{
		id:       "",
		name:     "",
		selector: &ps,
		poolIdcs: []*poolIdc{
			{
				idcId:   "fdsfsfdsfsdfdsrgtfdgfdgfdgdfgdfg343",
				groupId: "1",
				weight:  10,
			},
			{
				idcId:   "2",
				groupId: "4",
				weight:  10,
			},
		},
		p:        nil,
	}

	d := p.selectAddrs(&Context{idcAffinity: []string{"3"}},&Address{})
	fmt.Println(d)
	fmt.Println(ps.getBucket("fdsfsfdsfsdfdsrgtfdgfdgfdgdfgdfg343"))
}

func BenchmarkName(b *testing.B) {
	ps := &weightPoolSelector{

	}

	p:=&pool{
		id:       "",
		name:     "",
		selector: ps,
		poolIdcs: []*poolIdc{
			{
				idcId:   "2321435gfdgkfdgkfgfgdfdfsfds",
				groupId: "1",
				weight:  10,
			},
			{
				idcId:   "2321435gfdgkfdgkfgfgdfdfsfdb",
				groupId: "4",
				weight:  10,
			},
		},
		p:        nil,
	}
	ps.init(p)
	for i := 0; i < b.N; i++ {
		p.selectAddrs(&Context{idcAffinity: []string{"2321435gfdgkfdgkfgfgdfdfsfdg"}},&Address{})
		//ps.getBucket("fdsfsfdsfsdfdsrgtfdgfdgfdgdfgdfg343")
	}
	fmt.Println(ps.getBucket("fdsfsfdsfsdfdsrgtfdgfdgfdgdfgdfg343"))
}
