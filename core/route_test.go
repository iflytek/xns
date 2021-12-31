package core

import (
	"fmt"
	"testing"
)

func must(err error){
	if err != nil{
		panic(err)
	}
}

func TestRouteSelector_addRoute(t *testing.T) {
	rs := &RouteSelector{
		hosts: map[string]*ruleSelector{},
		genericHosts: nil,
		genericBack:  nil,
	}
	must(rs.addRoute([]string{"a.b.c"},[]*rule{{}},&route{priority: 2,}))
	must(rs.addRoute([]string{"a.b.c"},[]*rule{{args: []arg{newEqArg("a","1")}}},&route{priority:2,createAt: 4}))
	must(rs.addRoute([]string{"a.b.c"},[]*rule{},&route{priority: 1,}))
	must(rs.addRoute([]string{"*.a.b.c"},[]*rule{{}},&route{priority: 5,}))
	must(rs.addRoute([]string{"*.c.a.b.c"},[]*rule{{}},&route{priority: 6,}))

	fmt.Println(rs.getRoute(&Context{host: "a.b.a.b.c",params: map[string]string{
		"a":"1",
	}}))
}
